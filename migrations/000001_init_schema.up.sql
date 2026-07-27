CREATE TABLE customers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    telegram_user_id BIGINT NOT NULL UNIQUE,
    username TEXT,
    first_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE services (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL,
    price_minor_units INTEGER NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT service_duration
        CHECK (duration_minutes > 0),

    CONSTRAINT service_price
        CHECK (price_minor_units >= 0)
);

CREATE TABLE barbers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE barber_services (
    barber_id BIGINT NOT NULL
        REFERENCES barbers(id)
        ON DELETE CASCADE,

    service_id BIGINT NOT NULL
        REFERENCES services(id)
        ON DELETE CASCADE,

    PRIMARY KEY (barber_id, service_id)
);

CREATE TABLE barber_working_hours (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    barber_id BIGINT NOT NULL
        REFERENCES barbers(id)
        ON DELETE CASCADE,

    weekday SMALLINT NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,

    CONSTRAINT weekday_working
        CHECK (weekday BETWEEN 1 AND 7),

    CONSTRAINT time_valid
        CHECK (start_time < end_time)
);

CREATE TABLE bookings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    customer_id BIGINT NOT NULL
        REFERENCES customers(id),

    barber_id BIGINT NOT NULL
        REFERENCES barbers(id),

    service_id BIGINT NOT NULL
        REFERENCES services(id),

    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT booking_valid_time
        CHECK (starts_at < ends_at)
);