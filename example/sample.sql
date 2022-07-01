

CREATE TABLE IF NOT EXISTS doc (
      doc_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'unique identifier',
      ver INT NOT NULL COMMENT 'version, starts from zero',
      ref VARCHAR(50) NOT NULL COMMENT 'reference to locate the external file',
      deleted int NOT NULL DEFAULT '0' COMMENT '0:activated, others:deleted',
      create_user_id int NOT NULL COMMENT 'record create user ID',
      create_time datetime NOT NULL COMMENT 'record create time',
      update_user_id int NOT NULL COMMENT 'record update user ID',
      update_time datetime NOT NULL COMMENT 'record update time',
      recver int NOT NULL DEFAULT '0' COMMENT 'record version'
);
CREATE TABLE IF NOT EXISTS doc_tag (
      doc_tag_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'unique identifier',
      doc_id INT NOT NULL COMMENT 'doc id',
      tag_id INT NOT NULL COMMENT 'tag id',
      deleted int NOT NULL DEFAULT '0' COMMENT '0:activated, others:deleted',
      create_user_id int NOT NULL COMMENT 'record create user ID',
      create_time datetime NOT NULL COMMENT 'record create time',
      update_user_id int NOT NULL COMMENT 'record update user ID',
      update_time datetime NOT NULL COMMENT 'record update time',
      recver int NOT NULL DEFAULT '0' COMMENT 'record version'
);
CREATE TABLE IF NOT EXISTS tag (
      tag_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'unique identifier',
      code VARCHAR(10) COMMENT 'tag code. eg. premarket md application index A1',
      name VARCHAR(20) NOT NULL COMMENT 'tag name',
      priority VARCHAR(50) COMMENT 'tag priority',
      descr VARCHAR(50) COMMENT 'tag description',
      source CHAR(3) NOT NULL COMMENT 'Source; COM: common, PRE: pre, POS: post',
      type CHAR(3) NOT NULL COMMENT 'Tag Type; ATT: attachment, APP: application, COM: company',
      deleted int NOT NULL DEFAULT '0' COMMENT '0:activated, others:deleted',
      create_user_id int NOT NULL COMMENT 'record create user ID',
      create_time datetime NOT NULL COMMENT 'record create time',
      update_user_id int NOT NULL COMMENT 'record update user ID',
      update_time datetime NOT NULL COMMENT 'record update time',
      recver int NOT NULL DEFAULT '0' COMMENT 'record version'
);



ALTER TABLE IF EXISTS doc_tag ADD CONSTRAINT fk_doc_tag_doc_id FOREIGN KEY (doc_id) REFERENCES doc (doc_id);
ALTER TABLE IF EXISTS doc_tag ADD CONSTRAINT fk_doc_tag_tag_id FOREIGN KEY (tag_id) REFERENCES tag (tag_id);



DROP TABLE IF EXISTS doc;
DROP TABLE IF EXISTS doc_tag;
DROP TABLE IF EXISTS tag;



