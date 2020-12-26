<h1>Имплементация redis</h1>
<h2>Инструкция по запуску</h2>
<h3>Структура проекта</h3>
<a href="https://github.com/golang-standards/project-layout">https://github.com/golang-standards/project-layout</a>
<br>
<ul>
    <li>build - содержит все необходимое для запуска и работы сервера. В папке /server и /client лежат файлы для запуска сервера и клиента.</li>
    <li>cmd - содержит main package для сервера и клиента;</li>
    <li>deployments - содержит файл с docker-compose.yml. Контейнеры запускаются здесь;</li>
    <li>internal/server - содержит пакеты для работы сервера.
        <ul>
            <li>server - конфиг, сервер, роутер и обработчики;</li>
            <li>store - включает в себя реализацию интерфейса и сам интерфейс для работы с redis, mock для тестирования.</li>
        </ul>
    </li>
    <li>
    internal/client - реализация клиента
    </li>
    <li>
    test - json файлы для нагрузочных тестов
    </li>
</ul>
<h3>Запуск контейнеров</h3>
<p>В директории ./deployments прописать:
    <br>
    <code>docker-compose up -d</code>
    <br>
    После этого соберутся и запустятся контейнеры для работы с go и postgres
</p>
<p>
 Посмотреть их названия можно командой:
    <br>
    <code>
        docker ps
    </code>
</p>
<h3>Запуск сервера</h3>
<ul>
    <li>
        В директории ./build/server необходимо создать .env файл и заполнить его по примеру .env_example(скопировать все из .env_example в .env), аналогично для ./build/client
        <br>
        <code>
            cp .env_example .env
        </code>
    </li>
    <li>
        В директории ./build/server необходимо прописать команду для сборки сервера:
        <br>
            <code>go build ../../cmd/server/main.go</code>
        <br>
    </li>
    <li>
        В директории ./build/client необходимо прописать команду для сборки клиента:
        <br>
            <code>go build ../../cmd/lient/main.go</code>
        <br>
    </li>
    <li>
        Теперь их можно запустить командой:
        <br>
            <code>./main</code>
        <br>
    </li>
</ul>
<h3>
    Api методы для сервера
</h3>
<ul>
    <li>
        создание списка RPUSH
        <br>
        <code>
            curl -X POST
            -d '{"key":"asd"{"key":"list:1", "value":["ivan", 1, 3.2]}'
        </code>
        <br>
        результат
        <br>
        <code>
            http://127.0.0.1:3000/list/set
        </code>
    </li>
    <li>
            получение списка LRANGE
        <br>
        <code>
            curl -X GET
            http://127.0.0.1:3000/list/get?key=list:1
        </code>
        <br>
        результат
        <br>
        <code>
            {"error":"","result":["ivan",1,3.2]}
        </code>
    </li>
    <li>
            добавить строку SET
        <br>
        <code>
            curl -X POST
             -d '{"key":"user:1", "value":"string"}' http://127.0.0.1:3000/string/set
        </code>
        <br>
        результат
        <br>
        <code>
            {"error":"","result":"OK"}
        </code>
    </li>
        <li>
        получить строку по значению GET
        <br>
        <code>
            curl -X GET
            http://127.0.0.1:3000/string/get?key="user:1"  
        </code>
        <br>
        результат
        <br>
        <code>    
            {"error":"","result":"string"}
        </code>
=
    <li>
        получить значение поля hash HGET
        <br>
        <code>
        curl -X GET http://127.0.0.1:3000/hash/hget?key=user:Ivan&field=password
        </code>
        <br>
        результат
        <br>
        <code>
        {"error":"",result":"$2a$08$BN5DyPquIrPhAnTQNxtrEOAXxMZgPAzQdNYJydpgMXGuRBy6tRP76"}
        </code>
    </li>
    <li>
        получить значение поля hash HGetAll
        <br>
        <code>
        curl -X GET 127.0.0.1:3000/hash/get?key=user:Ivan  
        <code>
        <br>
        результат
        <br>
        </code>
        {"error":"","result":{"lastname":"Lapshin","login":"Ivan","password":"$2a$08$BN5DyPquIrPhAnTQNxtrEOAXxMZgPAzQdNYJydpgMXGuRBy6tRP76","role":"user"}}
        </code>
    </li>
    <li>
        список ключей KEYS
        <br>
        <code>
            curl -X GET http://127.0.0.1:3000/keys?pattern=*                        
        </code>
        <br>
        результат
        <br>
        <code>
            {"error":"","result":["Ivan","user:ivan","list:1","{\"name\":\"Ivan\",\"sex\":\"male\"}","user:Ivan","key","user:1","user:Ivan2","user:2","user:3"]}
        </code>
    </li>
    <li>
        удалить ключ DEL
        <br>
        <code>
        curl -X POST -d '{"key":"user:1"}' 127.0.0.1:3000/del
        </code>
    </li>
    <li>
        сохранить данные на диск SAVE
        <br>
        <code>
        curl -X POST  127.0.0.1:3000/save
        </code>
        <br>
        результат
        <br>
        <code>
        {"error":"","result":"Файл dump.rdb в папке [project name]/build/redis/data"}
        </code>
    </li>
</ul>