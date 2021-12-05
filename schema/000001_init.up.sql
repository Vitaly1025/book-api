CREATE TABLE genre(
    id serial not null unique,
    name varchar(255) not null,
    PRIMARY KEY(id)
);

CREATE TABLE book(
    id serial not null unique,
    amount integer not null,
    name varchar(255) not null unique,
    price decimal not null,
    genre int,
    PRIMARY KEY(id),
    CONSTRAINT fk_genre
        FOREIGN KEY(genre)
            REFERENCES genre(id)
);

INSERT INTO genre VALUES
    (1,'Adventure');
INSERT INTO genre VALUES
    (2,'Classics');
INSERT INTO genre VALUES
    (3,'Fantasy');
