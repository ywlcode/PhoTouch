DROP TABLE IF EXISTS user;  /* 用户信息 */
CREATE TABLE user (
  id         INT AUTO_INCREMENT NOT NULL,
  email     VARCHAR(255) NOT NULL,
  pwd       VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO user
  (email,pwd)
values('123@163.com','123'),
      ('123@qq.com','123');

select * from user;



DROP TABLE IF EXISTS imgpublic; /*  分享的图片 */
CREATE TABLE imgpublic (
  id         INT AUTO_INCREMENT NOT NULL,
  url        VARCHAR(255) NOT NULL,
  minurl     VARCHAR(255) NOT NULL,
  year       INT NOT NULL,
  month      INT NOT NULL,
  day        INT NOT NULL,
  userid     INT NOT NULL,
  PRIMARY KEY (`id`) 
);

DESCRIBE imgpublic;


DROP TABLE IF EXISTS imguser;  /*  私有的图片 */
CREATE TABLE imguser (
  id         INT AUTO_INCREMENT NOT NULL,
  url        VARCHAR(255) NOT NULL,
  minurl     VARCHAR(255) NOT NULL,
  year       INT NOT NULL,
  month      INT NOT NULL,
  day        INT NOT NULL,
  userid     INT NOT NULL,
  PRIMARY KEY (`id`) 
);

DESCRIBE imguser;

DROP TABLE IF EXISTS good;  /*  图片点赞数 */
CREATE TABLE good (
  userid INT NOT NULL,
  imgid  INT NOT NULL
);

DESCRIBE good;

