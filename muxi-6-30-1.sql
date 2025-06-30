create database MUXI;
use MUXI;

CREATE TABLE Book (
    ID VARCHAR(50) PRIMARY KEY,
    Title NVARCHAR(100) NOT NULL,
    Author NVARCHAR(50) NOT NULL
);

CREATE TABLE Storage (
    BookID VARCHAR(50) PRIMARY KEY,
    Stock INT NOT NULL DEFAULT 0,
    CONSTRAINT FK_Storage_Book FOREIGN KEY (BookID) REFERENCES Book(ID)
);

CREATE TABLE Person (
    ID INT PRIMARY KEY,
    Name NVARCHAR(50) NOT NULL
);

CREATE TABLE BorrowRecord (
    ID INT IDENTITY(1,1) PRIMARY KEY,
    PersonID INT NOT NULL,
    BookID VARCHAR(50) NOT NULL,
    BorrowDate DATETIME NOT NULL DEFAULT GETDATE(),
    ReturnDate DATETIME NULL,
    Status NVARCHAR(20) NOT NULL DEFAULT '借出',
    CONSTRAINT FK_BorrowRecord_Person FOREIGN KEY (PersonID) REFERENCES Person(ID),
    CONSTRAINT FK_BorrowRecord_Book FOREIGN KEY (BookID) REFERENCES Book(ID),
    CONSTRAINT CHK_Status CHECK (Status IN ('借出', '已归还', '逾期'))
);

-- 插入书籍信息
INSERT INTO Book (ID, Title, Author) VALUES
('go-away', 'the way to go', 'Ivo'),
('go-lang', 'Go语言圣经', 'Alan,Brian'),
('go-web', 'Go Web编程', 'Anonymous'),
('con-cur', 'Concurrency in Go', 'Katherine');

-- 插入库存信息
INSERT INTO Storage (BookID, Stock) VALUES
('go-away', 20),
('go-lang', 17),
('go-web', 4),
('con-cur', 9);

-- 插入人员信息
INSERT INTO Person (ID, Name) VALUES
(1, '小明'),
(2, '张三'),
(3, '翟曙');

-- 插入喜好关系
INSERT INTO BorrowRecord (PersonID, BookID, Status) VALUES
(1, 'go-away', '借出'),
(1, 'go-web', '借出'),
(1, 'con-cur', '借出'),
(2, 'go-away', '借出'),
(3, 'go-web', '借出'),
(3, 'con-cur', '借出');

SELECT p.Name 
FROM Person p
JOIN BorrowRecord br ON p.ID = br.PersonID
WHERE br.BookID = 'go-web';

SELECT b.ID, b.Author, b.Title, s.Stock
FROM Book b
JOIN Storage s ON b.ID = s.BookID
WHERE b.ID NOT IN (SELECT DISTINCT BookID FROM BorrowRecord);

SELECT p.Name AS 人名, b.Title AS 书名
FROM Person p
JOIN BorrowRecord br ON p.ID = br.PersonID
JOIN Book b ON br.BookID = b.ID
ORDER BY p.Name, b.Title;