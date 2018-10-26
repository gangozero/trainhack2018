CREATE EXTENSION pgcrypto;

CREATE TABLE orders
(
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    train TEXT NOT NULL,
    carriage TEXT NOT NULL,
    station TEXT NOT NULL,
    repeat_order BOOLEAN NOT NULL DEFAULT false,
    delivery BOOLEAN NOT NULL DEFAULT false,
    ord JSONB NOT NULL DEFAULT '[]'::jsonb,
    ts_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ts_ready TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
