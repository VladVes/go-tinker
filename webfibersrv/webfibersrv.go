package webfibersrv

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/VladVes/go-tinker/v2/data"
	"github.com/VladVes/go-tinker/v2/data/entities"
	"github.com/VladVes/go-tinker/v2/data/schemas"
	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const profileUnknown = "unknown"

func Run() {
	app := fiber.New()

	// -------------------------------------Middlewares and Route Groups-----------------------------------------
	// Пример простейшей middleware (определены в middlewares.go) и её использования
	// будет запускатся перед обработкой каждого запроса вне зависимости от пути
	app.Use(TestMiddleware1)

	app.Get("/tests", func(c *fiber.Ctx) error {
		return c.SendString("Middleware test!\n")
	})

	// Группировка маршрутов что бы связать заданную middleware только с ними
	// Создаем группу с префиксом пути запроса "/mw_tests".
	midTestGroup := app.Group("/mw_tests")
	// устанавливаем middleware для группы
	midTestGroup.Use(TestMiddleware2)
	midTestGroup.Get("/aciton", func(c *fiber.Ctx) error {
		return c.SendString("route with middleware!")
	})
	// если для группы не передавать явно middleware то никакая не зпустится либо
	// запустится только общая для всех маршрутов
	midTestGroup2 := app.Group("/no_mw_tests")
	midTestGroup2.Get("/action", func(c *fiber.Ctx) error {
		return c.SendString("route with no middleware")
	})

	// Пример подключения готовой мидлвары - логгера запросов
	// --
	// хорошей практикой считается логировать идентификатор, который поможет нам связать все
	// логи в рамках одного запроса. Для этого нам нужно подключить к проекту
	// пакет github.com/gofiber/fiber/v2/middleware/requestid.
	// И инициализировать его перед посредником для логирования
	app.Use(requestid.New())
	// По умолчанию логирование происходит в консоль, но можно настроить логирование в файл или в другие системы логирования.
	// при инициализации посредника для логирования мы можем указать формат логов, который подходит для нашего проекта
	app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid}: ${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000000",
	}))
	app.Get("/logger_test", func(c *fiber.Ctx) error {
		time.Sleep(300 * time.Millisecond)

		logrus.WithFields(logrus.Fields{
			"request_id": c.Locals("requestid"),
		}).Warn("something went wrong")

		return c.SendString("OK")
	})
	// Чтобы защититься со стороны веб-приложения от атак, следует настроить ограничение количества запросов — throttling.
	// Для этого мы будем использовать пакет github.com/gofiber/fiber/v2/middleware/limiter.
	// пример использования миддлвары ограничивающей число запросов:
	app.Use(limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Max:        3,
		Expiration: 10 * time.Second,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// -------------------------------------First route and response-----------------------------------------
	app.Get("/address", func(c *fiber.Ctx) error {
		fmt.Println("Processing!")
		return c.SendString("Hello, go! Don't gitve up! Fight!\n")
	})

	// -------------------------------------Query---------------------------------------------------------------
	// работа с параметрами запроса
	app.Get("/profile", func(c *fiber.Ctx) error {
		profileId := c.Query("profile_id", profileUnknown)
		if profileId == "" {
			profileId = profileUnknown
		}
		if profileId == profileUnknown {
			return c.Status(fiber.StatusUnprocessableEntity).SendString("porofile_id is required")
		}

		return c.SendString(fmt.Sprintf("user Profile ID is: %s", profileId))
	})

	// ------------------------------------Dynamin route----------------------------------------------------------------
	// пример динамического роутинга
	app.Get("/likes/:postId", func(c *fiber.Ctx) error {
		postId := c.Params("postId")
		likes, ok := data.PostLikes[postId]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"postId": postId,
			}).Info(fmt.Sprintf("A post with this id is not found: %s", postId))
			return c.Status(fiber.StatusNotFound).SendString("no post Id\n")
		}

		return c.SendString(fmt.Sprintf("Post id: %s, Likes: %s\n", postId, strconv.FormatInt(likes, 10)))
	})

	app.Post("/likes/:postId", func(c *fiber.Ctx) error {
		postId := c.Params("postId")

		data.PostLikes[postId]++
		likes := data.PostLikes[postId]

		status := fiber.StatusOK
		if likes == 1 {
			status = fiber.StatusCreated
		}
		return c.Status(status).SendString(fmt.Sprintf("Post id: %s, Likes: %s\n", postId, strconv.FormatInt(likes, 10)))

	})

	// -------------------------------JSON-----------------------------------------------------------------
	// Пример десериализации данных передаваемых в теле запсроса
	// и сериализации отправляемого ответа
	app.Post("/logs", func(c *fiber.Ctx) error {
		var request schemas.CreateLogEntryRequest
		if err := c.BodyParser(&request); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		logEntry := schemas.LogEntry{
			ID:        uuid.New().String(),
			Message:   request.Message,
			Level:     request.Level,
			Timestamp: request.Timestamp,
		}

		// Упрощенное хранение в памяти приложения
		data.Logs = append(data.Logs, logEntry)

		return c.JSON(schemas.CreateLogEntryResponse{
			ID: logEntry.ID,
		})
	})

	// ----------------------------------Logrus POST JSON---------------------------------------------------------
	// Пример обработки POST запроса в теле которого передаётся отсортированый массив чисел
	// и искомое число, в ответе в теле возвращается JSON с полем target_index с индексом искомого числа в массиве
	// В примере используются две стркутруты - схемы для сериализации и десериализации. Схемы имеют тегированные поля
	// в которых содержится описание того из каких полей JSON что использовать для каких полей структуры и на оборот
	//  какой JSON должен получиться из каких полей
	app.Post("/searchIndex", func(c *fiber.Ctx) error {
		req := schemas.SearchIndexReq{}
		resp := schemas.SearchIndexResp{
			TargetIndex: -1,
		}

		if err := c.BodyParser(&req); err != nil {
			resp.Error = "Invalid JSON"
			logrus.WithError(err).Error(resp.Error)
			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		targetIndex := slices.Index(req.Numbers, req.Target)
		if targetIndex == -1 {
			resp.Error = "Target not found"
			logrus.WithFields(logrus.Fields{
				"target":  req.Target,
				"numbers": req.Numbers,
			}).Info(resp.Error)
			return c.Status(fiber.StatusNotFound).JSON(resp)
		}

		resp.TargetIndex = targetIndex
		return c.Status(fiber.StatusOK).JSON(resp)
	})

	// -----------------------------------Layered arch-----------------------------------------------------------------
	// Пример простейшей логики создания и получения заказа реализованный по принципу слоёной архитекутуры
	// Слой обработчика -> слой бизнеc-логики (тут опущен) -> слой хранилица
	// определены в пакете entities
	orderHandler := &entities.OrderHandler{
		Storage: &entities.OrderStorage{
			Orders: data.Orders,
		},
	}
	app.Post("/orders", orderHandler.CreateOrder)
	app.Get("/orders/:id", orderHandler.GetOrder)

	// -------------------------------------CRUD------------------------------------------------------------
	// Пример CRUD по сущности Employee c хранением в памяти
	//
	employeeHandler := &entities.EmployeeHanlder{
		Storage: &entities.EmployeeStorageInMemory{
			Employee: data.Employee,
		},
	}

	app.Post("/employees", employeeHandler.CreateEmployee)
	app.Get("/employees", employeeHandler.GetEmployees)
	app.Get("/employees/:id", employeeHandler.GetEmployeeByID)
	app.Patch("/employees/:id", employeeHandler.UpdateEmployee)
	app.Delete("/employees/:id", employeeHandler.DeleteEmployee)

	// ------------------------------------Validation go-playgorund/validator---------------------------------
	// Пример простой ручной валидации запроса по созданию поста

	app.Post("/post", func(c *fiber.Ctx) error {
		// CreatePostReq определена в Validators.go
		var req CreatePostReq
		if err := c.BodyParser(&req); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		// метод определен в Validators.go, проверяет поля структруры и возвращает ошибку c определенным текстом если поле не соответсвтует
		err := req.Validate()
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}
		// @TODO сохранение данных в бд

		return c.SendStatus(fiber.StatusOK)
		// curl --location --request POST 'http://localhost/post' \
		// --header 'Content-Type: application/json' \
		// --data-raw '{"user_id": -1, "text": ""}'
	})

	// Пример использвоания библиотеки go-playground/validator https://pkg.go.dev/github.com/go-playground/validator/v10
	// валидирующей поля используя их аннотации
	type CreatePostRequest struct {
		// Описываем правила валидации в аннотациях полей структуры.
		UserID int64  `json:"user_id" validate:"required,min=0"`
		Text   string `json:"text" validate:"required,max=140"`
		// имеет множество различных готовых функций для проверки, к примреу проверка корректности email:
		Email string `json:"email" validate:"required,email"`
	}

	validate := validator.New()

	app.Post("/posts", func(c *fiber.Ctx) error {
		var req CreatePostRequest
		if err := c.BodyParser(&req); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		// Проверка запроса
		err := validate.Struct(req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		// @TODO создание поста и запись в бд

		return c.SendStatus(fiber.StatusOK)
		// curl --location --request POST 'http://localhost/posts' \
		// --header 'Content-Type: application/json' \
		// --data-raw '{"user_id": -1, "text": ""}'

	})
	// Пользователься валидация в go-playground/validator
	// Например, мы хотим проверить, что в публикуемом посте отсутствуют слова-фильтры.
	type CreateNewPostRequest struct {
		// Описываем правила валидации в аннотациях полей структуры.
		// allowable_text - специальная аннотация
		UserID int64  `json:"user_id" validate:"required,min=0"`
		Text   string `json:"text" validate:"required,max=140,allowable_text"`
	}
	// запрещенные слова
	var forbiddenWords = []string{
		"umbrella",
		"shinra",
	}

	validateWithAllowable := validator.New()

	vErr := validateWithAllowable.RegisterValidation("allowable_text", func(fl validator.FieldLevel) bool {
		text := fl.Field().String()
		for _, word := range forbiddenWords {
			if strings.Contains(strings.ToLower(text), word) {
				return false
			}
		}
		return true
	})
	if vErr != nil {
		log.Fatal("register validation ", vErr)
	}

	app.Post("/newposts", func(ctx *fiber.Ctx) error {
		// Парсинг JSON-строки из тела запроса в объект.
		var req CreateNewPostRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		// Проверка запроса на корректность.
		err := validateWithAllowable.Struct(req)
		if err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		// @TODO Сохранение поста в хранилище.

		return ctx.SendStatus(fiber.StatusOK)
		// curl --location --request POST 'http://localhost:8080/newposts' \
		// --header 'Content-Type: application/json' \
		// --data-raw '{"user_id": 100, "text": "Hello Umbrella corp!"}'
	})

	//******************************JWT USER AUTH********************************
	authHandler := &entities.AuthHandler{
		Storage: &entities.AuthStorage{
			Users: data.Users,
		}}

	publicGroup := app.Group("")

	publicGroup.Post("/register", authHandler.Register)
	publicGroup.Post("/login", authHandler.Login)

	authorizedGroup := app.Group("/user")

	// потому как data.Users это map который является ссылочным типом
	// то оба экземпляра и authHandler и authHandler будут внутри себя
	// использовать один и тот же мап (мап при передаче или присваивании
	// копирует тольк ссылку на структуру )
	userHandler := &entities.UserHandler{
		Storage: &entities.AuthStorage{
			Users: data.Users,
		},
	}
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: entities.JwtSecretKey,
		},
		ContextKey: entities.ContextKeyUser,
	}))
	authorizedGroup.Get("/profile", userHandler.Profile)

	// регистрация
	// 	curl --location --request POST 'http://localhost/register' \
	// --header 'Content-Type: application/json' \
	// --data-raw '{
	//     "email": "john@doe.com",
	//     "name": "John",
	//     "password": "pickles"
	// }'

	// вход
	// 	curl --location --request POST 'http://localhost:8080/login' \
	// --header 'Content-Type: application/json' \
	// --data-raw '{
	//     "email": "john@doe.com",
	//     "password": "pickles"
	// }'

	// получение данных пользователя
	// 	curl -v 'http://localhost:8080/user/profile' \
	// --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njg5NTE0NDEsInN1YiI6ImpvaG5AZG9lLmNvbSJ9.e4yIoGzQC8ckcRISBjt4g18S2VEBiHrRhXG7N39-7qI'

	// ----------------------------------------------------------------------------------------------------
	logrus.Fatal(app.Listen(":" + HttpPort))
}
