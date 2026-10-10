package basics

import (
	"fmt"
	"time"
)

// 38. Как создать метод и свзяать его со структурой User (Name и Age) так что бы выводить приветсвтие используя поля этой структуры?
// Где должен быть определён (описан) этот метод?
// Как вызывать этот метод?
// Сделать пример со структурой User имеющей поля имя и возраст и методом который выводит эти поля в формате Имя: Возраст:
type User struct {
	Name  string
	Email string
	Age   int
}

func (u User) PrintUser() {
	fmt.Printf("Имя: %s\nВозраст: %d\n", u.Name, u.Age)
}

// 48. Как возвращаются структруры из функций?
// Как создать функцию NewUser которая получает имя и возраст
// создаёт структуру type User struct { Name string Age int }
// и возвращает указатель на неё?
// Зачем возвращать указатели на структуры?
func NewUser(name string, age int) *User {
	u := User{
		Name: name,
		Age:  age,
	}
	// - Возможность изменять поля структуры вне функции.
	// - Экономию памяти (особенно при работе с большими структурами).
	// - Консистентность в стиле, если структура часто используется по ссылке.
	return &u
}

// 53. Cоздать структуру User с полем Name и методы для неё Log (ничего не делает с данными стуркруры) и PrintName (выводящий поля Name). В main создать указатель usrPtr на User но без явного определения т.е. так: var usrPtr *User .
// Что будет если попытаться вызвать метод Log?
// И что будет если попытаться вызвать метод PrintName?
// При каком изменении метода Log результат обработки вызова этого метода измениться и почему?
func (u *User) Log() {
	fmt.Println("some log, no user field")
}

// Можно такой методо определить и как статический т.е. без определения экземплякар структуры:
func (*User) Log2() {
	fmt.Println("some log, no user field")
}

func (u *User) PrintUserName() {
	fmt.Println(u.Name)
}

func UserProcess() {
	var u *User
	u.Log() // ошибки не будет даже учитывая что u указывает в никуда - nil
	// u.PrintName() // - panic: invalid memory address or nil pointer dereference
	u.Log2() // тоже работает
}

// 39.1 Как можно определить метод который будет вести себя так же как static методы в других ЯП?
// Т.е. будет делать что-то не связанное с конкретным экземпляром структуры. Как можно вызвать такой метод?

func (User) ModelDesc() string {
	return "User model meta info"
}

// 156. Какой есть механизм у структур позволяющий реализовывать композицию?

type EmployeePayload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type GetEmployeeResponse struct {
	// Анонимное поле позволяет встраивать поля одной структуры в другую - это механизм для реализации композиции.
	// Поле не имеет имени — вместо этого указывается только тип.
	// Все поля и методы  встраиваемой структуры EmployeePayload становятся доступными напрямую
	// через экземпляр внешней структуры (т.е. так getEmployeeResponse.Email или getEmployeeResponse.Role ).
	EmployeePayload
	// Приоритет полей: если у GetEmployeeResponse и EmployeePayload
	// есть поля с одинаковыми именами, приоритет имеет поле внешней структуры.
	// Гибкость: позволяет переиспользовать поля и методы без дублирования кода.
}

type ListEmployeesResponse struct {
	Employees []GetEmployeeResponse `json:"employees"`
}

// 157. Как в Go происходит приведение составных типов а именно одной структуры к другой что нужно при этом помнить?
// Предположим, в хранилище используется структура:
type Employee struct {
	ID        string
	Email     string
	Role      string
	CreatedAt time.Time
}

// а для создания ответа используется структура:
type EmployeeRespPayload struct {
	ID    string
	Email string
	Role  string
}

// Нужно привести Employee для отправки http ответа к EmployeeRespPayload
// автоматического механизма для такого приведение нет, нужно определять поля руками
func StructTypeConversionDemo() EmployeeRespPayload {
	emp := Employee{
		ID:    "123asdf",
		Email: "empl@mail.com",
		Role:  "employee",
	}

	return EmployeeRespPayload{
		ID:    emp.ID,
		Email: emp.Email,
		Role:  emp.Role,
	}
}
