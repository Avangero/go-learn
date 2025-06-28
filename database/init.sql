-- Создание баз данных для каждого микросервиса
CREATE DATABASE auth_db;
CREATE DATABASE user_db;  
CREATE DATABASE course_db;

-- Подключаемся к auth_db для создания таблиц
\c auth_db;

-- Создание таблицы пользователей для Auth Service
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('employee', 'manager')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Создание индексов для быстрого поиска
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Подключаемся к user_db для будущих таблиц
\c user_db;

-- Создание таблицы профилей пользователей для User Service
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, -- FK to auth.users
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    position VARCHAR(100),
    department VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индекс для связи с auth.users
CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);

-- Подключаемся к course_db для будущих таблиц
\c course_db;

-- Создание таблиц для Course Service
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    skills JSONB, -- массив навыков в JSON формате
    created_by UUID NOT NULL, -- FK to auth.users (manager)
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE course_modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    order_num INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, -- FK to auth.users
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    completed_modules JSONB DEFAULT '[]', -- массив UUID завершенных модулей
    status VARCHAR(50) NOT NULL CHECK (status IN ('enrolled', 'in_progress', 'completed')),
    progress DECIMAL(5,4) DEFAULT 0.0000, -- прогресс от 0.0000 до 1.0000
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Индексы для Course Service
CREATE INDEX idx_courses_created_by ON courses(created_by);
CREATE INDEX idx_course_modules_course_id ON course_modules(course_id);
CREATE INDEX idx_course_modules_order ON course_modules(course_id, order_num);
CREATE INDEX idx_user_progress_user_id ON user_progress(user_id);
CREATE INDEX idx_user_progress_course_id ON user_progress(course_id);
CREATE INDEX idx_user_progress_status ON user_progress(status); 