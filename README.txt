1. Установка 

go get -u github.com/DimkaTheGreat/sittme
go build

2. Запуск

./sittme -p <порт> -timeout=<таймаут в секундах>

3. Роутинг

GET http://<ip-address:port/list - отображение списка текущих трансляций и их состояние
GET http://<ip-address:port/create -  создание уникального идентификатора трансляции
DELETE http://<ip-address:port/delete?id=<id трансляции> - удаление трансляции по идентификатору
GET http://<ip-address:port/activate?id=<id трансляции> - команда запуска: установка состояния “Active”
GET http://<ip-address:port/interrupt?id=<id трансляции> - команда прерывания: установка состояния “Interrupted”, запуск таймера для перехода в состояние “Finished”
