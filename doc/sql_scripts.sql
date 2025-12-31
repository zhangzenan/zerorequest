
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
