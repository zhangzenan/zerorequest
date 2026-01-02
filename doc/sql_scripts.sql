
CREATE TABLE focus (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    pic VARCHAR(255) NOT NULL,
    link VARCHAR(255),
    position INT DEFAULT 0,
    status INT DEFAULT 1
);

INSERT INTO focus (id, title, pic, link, position, status) VALUES
('1', '首页轮播图1', 'https://example.com/images/banner1.jpg', 'https://example.com/link1', 1, 1),
('2', '首页轮播图2', 'https://example.com/images/banner2.jpg', 'https://example.com/link2', 2, 1),
('3', '首页轮播图3', 'https://example.com/images/banner3.jpg', 'https://example.com/link3', 3, 1),
('4', '首页轮播图4', 'https://example.com/images/banner4.jpg', 'https://example.com/link4', 4, 1),
('5', '首页轮播图5', 'https://example.com/images/banner5.jpg', 'https://example.com/link5', 5, 1);



CREATE TABLE `user` (
                        `id` varchar(255) NOT NULL COMMENT '用户ID',
                        `name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
                        `age` int(11) NOT NULL DEFAULT 0 COMMENT '年龄',
                        `sex` int(11) NOT NULL DEFAULT 0 COMMENT '性别(1男,2女)',
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';


INSERT INTO `user` (`id`, `name`, `age`, `sex`) VALUES
                                                    ('1', '张三', 25, 1),
                                                    ('2', '李四', 28, 2),
                                                    ('3', '王五', 30, 1),
                                                    ('4', '赵六', 22, 2),
                                                    ('5', '钱七', 26, 1);
