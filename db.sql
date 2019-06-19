CREATE TABLE "user" (
"id"  SERIAL ,
"login" VARCHAR(50) NOT NULL UNIQUE ,
"password" VARCHAR(60) NOT NULL ,
PRIMARY KEY ("id")
);

CREATE TABLE "video" (
"id"  SERIAL ,
PRIMARY KEY ("id")
);

CREATE TABLE "search_history" (
"id"  SERIAL ,
"query" VARCHAR(50) NOT NULL ,
"time" DATETIME NOT NULL ,
"id_user" INTEGER ,
PRIMARY KEY ("id")
);

CREATE TABLE "user_likes" (
"id"  SERIAL ,
"id_user" INTEGER ,
"id_video" INTEGER ,
PRIMARY KEY ("id")
);

ALTER TABLE "search_history" ADD FOREIGN KEY ("id_user") REFERENCES "user" ("id");
ALTER TABLE "user_likes" ADD FOREIGN KEY ("id_user") REFERENCES "user" ("id");
ALTER TABLE "user_likes" ADD FOREIGN KEY ("id_video") REFERENCES "video" ("id");
