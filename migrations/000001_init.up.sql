CREATE TABLE genre(
    id serial not null unique,
    name varchar(255) not null unique,
    PRIMARY KEY(id)
);

CREATE TABLE author(
    id serial not null unique,
    name varchar(255) not null unique,
    PRIMARY KEY(id)
);

CREATE TABLE book(
    id serial not null unique,
    pagecount integer not null,
    rate smallint,
    cover bytea,
    name varchar(255) not null unique,
    description text not null,
    authorId int,
    PRIMARY KEY(id),
    CONSTRAINT fk_author
        FOREIGN KEY(authorId)
            REFERENCES author(id)
);

CREATE TABLE book_genre(
    book_id int not null,
    genre_id int not null,
    PRIMARY KEY (book_id,genre_id),
    CONSTRAINT fk_book
        FOREIGN KEY (book_id) REFERENCES book(id)
    CONSTRAINT fk_genre
        FOREIGN KEY (genre_id) REFERENCES genre(id),
);