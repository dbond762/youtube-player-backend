CREATE TABLE "user" (
"id"  SERIAL ,
"login" VARCHAR(50) NOT NULL UNIQUE ,
"password" VARCHAR(60) NOT NULL ,
PRIMARY KEY ("id")
);

CREATE TABLE "video" (
"id"  VARCHAR(20) NOT NULL,
"title" VARCHAR(50) NOT NULL,
"pub_date" TIMESTAMP NOT NULL,
"description" TEXT NOT NULL,
"thumbnail" VARCHAR(70) NOT NULL,
"player" TEXT NOT NULL DEFAULT '',
PRIMARY KEY ("id")
);

CREATE TABLE "search_history" (
"id"  SERIAL ,
"query" VARCHAR(50) NOT NULL ,
"time" TIMESTAMP NOT NULL ,
"id_user" INTEGER ,
PRIMARY KEY ("id")
);

CREATE TABLE "user_likes" (
"id_user" INTEGER ,
"id_video" VARCHAR(20) ,
PRIMARY KEY ("id_user", "id_video")
);

ALTER TABLE "search_history" ADD FOREIGN KEY ("id_user") REFERENCES "user" ("id");
ALTER TABLE "user_likes" ADD FOREIGN KEY ("id_user") REFERENCES "user" ("id");
ALTER TABLE "user_likes" ADD FOREIGN KEY ("id_video") REFERENCES "video" ("id");
