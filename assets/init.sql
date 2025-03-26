-- Таблица пользователей, содержащая данные о зарегистрированных пользователях (почта, логин, пароль, аватар и др.)
CREATE TABLE IF NOT EXISTS "users" (
                                       "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
                                       "email" text UNIQUE NOT NULL,
                                       "username" text NOT NULL,
                                       "city" text NOT NULL DEFAULT 'Москва',
    "password" text NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               PRIMARY KEY ("id")
    );
