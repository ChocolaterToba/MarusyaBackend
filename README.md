# Конструктор скиллов El Marusia

## Как запустить
1. Установить Docker 
   ``` https://docs.docker.com/engine/install/ ```
2. В папке server/settings/values_local.yaml, изменить текст команд на необходимый вам.
3. Если вы хотите использовать базу в докере, то не меняйте ничего в настройках БД
4. Запустите контейнер ``` sudo docker-compose -f docker-compose.yaml up --build -d ```
5. Отлично. Вы запустили контейнер, подключиться к марусе можно через localhost:8080/api/marusia

## Как добавить свой скилл
1. Запустить сервер по инструкции выше
2. Создать отдельный файл для каждого теста (сценария), пример можно найти в 
   > server/settings/tests_example.csv
3. Отправить PUT Form-Data запрос на ```localhost:8080/api/test/add ```
   С полями: quizAmount: 1, 2, 3, etc; file1; file2; file3 where fileN
   (ждем апдейта с упрощением этой процедуры)
4. Можно тестировать.


## Отладка скилла через skill-debugger
1. Подключаемся к 
   ``` https://skill-debugger.marusia.mail.ru/ ```
2. Указываем адрес вашей машины в поле > Webhook url
   Например ``` my.marusia_server.ru:8080/api/marusia ```
3. [Опционально] Используем веб версию дебагера
4. [Опционально] В верхнем правом углу жмем > Подключение клиента Маруси 
   И следуем инструкциям

## Регистрация скилла 
```https://dev.vk.com/marusia/getting-started```

# MarusyaBackend

How do i launch this???
1. Install Docker
2. Go to server/settings/values_local.yaml, change special texts for ones you need
3. If you use dockerized db, don't change any database configs
4. $ sudo docker-compose -f docker-compose.yaml up --build -d
5. Done! Now your Marusia can connect to localhost:8080/api/marusia to get access to skills

How do i add skills???
1. Launch server per instructions above
2. Create separate files for each test, for example see server/settings/tests_example.csv
3. Send PUT Form-Data request to localhost:8080/api/test/add. Fields: quizAmount: 1, 2, 3, etc; file1; file2; file3 where fileN - files created during previous step
4. Done!
