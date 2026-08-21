-- public.appa_users definition

-- Drop table

-- DROP TABLE public.appa_users;
CREATE TABLE public.appa_users
(
    id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
    firebase_uid character varying(128) NOT NULL,
    email character varying(128) NOT NULL,
    is_admin boolean DEFAULT false,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    CONSTRAINT appa_users_pkey PRIMARY KEY (id),
    CONSTRAINT appa_users_firebase_uid_key UNIQUE (firebase_uid),
    CONSTRAINT appa_users_email_key UNIQUE (email)
);
CREATE INDEX idx_appa_users_firebase_uid ON public.appa_users USING btree (firebase_uid);
CREATE INDEX idx_appa_users_email ON public.appa_users USING btree (email);

-- public.payments_methods definition
-- Drop table
-- DROP TABLE public.payments_methods;
CREATE TABLE public.payments_methods
(
    id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
    name character varying(64) NOT NULL UNIQUE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    CONSTRAINT payments_methods_pkey PRIMARY KEY (id),
    CONSTRAINT payments_methods_name_key UNIQUE (name)
);

-- public.appa_manual_orders definition
-- Drop table
-- DROP TABLE public.appa_manual_orders;
CREATE TABLE public.appa_manual_orders
(
    id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
    order_name character varying(128) NOT NULL UNIQUE,
    order_id bigint NOT NULL,
    payment_method_id int4 NOT NULL,
    bill_image_url character varying(128) NOT NULL,
    amount decimal(10,2) NOT NULL,
    order_total_amount decimal(10,2) NOT NULL,
    requires_change bool NOT NULL,
    validate_status character varying(32) NOT NULL,
    return_data jsonb,
    logistic_validate bool DEFAULT false,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    CONSTRAINT appa_manual_orders_pkey PRIMARY KEY (id),
    CONSTRAINT appa_manual_orders_order_name_key UNIQUE (order_name)
);
CREATE INDEX idx_appa_manual_orders_order_id ON public.appa_manual_orders USING btree (order_id);
CREATE INDEX idx_appa_manual_orders_payment_method_id ON public.appa_manual_orders USING btree (payment_method_id);

ALTER TABLE public.appa_manual_orders ADD CONSTRAINT appa_manual_orders_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payments_methods (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE;
