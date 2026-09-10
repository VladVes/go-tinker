-- Запрос для генерации тестовых данных
INSERT INTO app_movies (title, genre, released_at, description, rating, votes)
SELECT
    'Movie ' || gs AS title,

    -- Случайный жанр из списка
    (ARRAY[
        'Action', 'Comedy', 'Drama', 'Horror', 'Sci-Fi',
        'Romance', 'Thriller', 'Fantasy', 'Animation', 'Documentary'
    ])[floor(random() * 10 + 1)] AS genre,

    -- Случайная дата с 1980 по 2023 год
    DATE '1980-01-01' + (random() * (DATE '2023-12-31' - DATE '1980-01-01'))::int AS released_at,

    -- Псевдо-описание
    'Description for movie ' || gs || '. ' ||
    substr(md5(random()::text), 1, 20) AS description,

    -- Рейтинг от 1.0 до 9.9
    round((random() * 8.9 + 1)::numeric, 1) AS rating,

    -- Количество голосов до 1 000 000
    (random() * 1000000)::bigint AS votes

FROM generate_series(1, 100) AS gs;