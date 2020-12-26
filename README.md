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
        Теперь их можно запустить командой (в /build/server и /build/client):
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
        </code>
        <br>
        результат
        <br>
        <code>
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
<h3>Api методы для клиента</h3>
<ul>
    <li>
            ддобавить список
        <br>
        <code>
        curl -X POST -d '{"method":"set", "payload": {"key": "client", "value":["1","2","3"], "ttl":10} }' http://127.0.0.1:3001/list
        </code>
        <br>
            результат
        <br>
        <code>
        {"error":"","result":{"error":"","result":"success"}}
        </code>
    </li>
    <li>
            получить список
        <br>
        <code>
        curl -X POST -d '{"method":"get", "payload": {"key": "client"} }' http://127.0.0.1:3001/list
        </code>
        <br>
            результат
        <br>
        <code>
        {"error":"","result":{"error":"","result":["1","2","3"]}}
        </code>
    </li>
    <li>
            создать строку
        <br>
        <code>
curl -X POST -d '{"method":"set", "payload": {"key": "client:2", "value":"string", "ttl":10} }' http://127.0.0.1:3001/string
        </code>
        <br>
            результат
        <br>
        <code>
{"error":"","result":{"error":"","result":"OK"}}
        </code>
    </li>
    <li>
            получить строку
        <br>
        <code>
curl -X POST -d '{"method":"get", "payload": {"key": "client:2"} }' http://127.0.0.1:3001/string
        </code>
        <br>
            результат
        <br>
        <code>
{"error":"","result":{"error":"","result":"string"}}
        </code>
    </li>
    <li>
    создать hash
    <br>
    <code>
            curl -X POST -d '{"method":"set", "payload": {"key": "client:3", "value":{"name": "Ivan", "age":10}, "ttl":10} }' http://127.0.0.1:3001/map
    </code>
    <br>
            результат
    <br>
    <code>
            {"error":"","result":{"error":"","result":"success"}}
    </code>
    </li>
    <li>
        получить hash
        <br>
        <code>
curl -X POST -d '{"method":"get", "payload": {"key": "client:3"} }' http://127.0.0.1:3001/map
        </code>
        <br>
            результат
        <br>
        <code>
{"error":"","result":{"error":"","result":{"age":"10","name":"Ivan"}}}
        </code>
    </li>
    <li>
        удалить ключ
        <br>
        <code>
curl -X POST -d '{"key": "client"}' http://127.0.0.1:3001/delete
        </code>
        <br>
            результат
        <br>
        <code>
{"error":"","result":{"error":"","result":1}}
        </code>
    </li>
    <li>
        получить список ключей
        <br>
        <code>
curl -X GET http://127.0.0.1:3001/keys?pattern=*      
        </code>
        <br>
            результат
        <br>
        <code>
{"error":"","result":{"error":"","result":["Ivan","user:ivan","list:1","{\"name\":\"Ivan\",\"sex\":\"male\"}","user:Ivan","key","user:Ivan2","user:2","client:3","client:2","user:3"]}}
        </code>
    </li>
</ul>
<h3>
Авторизация
</h3>
<ul>
    <li>
    Регистрация
    <br>
    <code>
    http://127.0.0.1:3000/signup
    </code>
    <br>
    Тело запроса
    <pre>
    {
        "login": "Ivan",
        "password": "asd"
    }
    </pre>
    </li>
    <li>
    Авторизация
    <br>
    <code>
    http://127.0.0.1:3000/login
    </code>
    <br>
    Тело запроса
    <pre>
    {
        "login": "Ivan",
        "password": "asd"
    }
    </pre>
    </li>
        <li>
    выход
    <br>
    <code>
    http://127.0.0.1:3000/logout
    </code>
    </li>
</ul>

<h3>
 Нагрузочное тестирование
</h3>
Чтобы начать тестирование необходимо зайти в папку /test и выполнить команду. Для описания флагов можно обратиться к Репозиторию с <a src="https://github.com/a696385/go-meter">модулем</a>
<br>
<code>
go-meter -t 12 -c 400 -d 30s -u http://127.0.0.1:3000/hash/set -s hash_set.json -v -m POST
</code>
<br>
результат
<br>
<pre>
Stats:            Min       Avg       Max
  Latency          0s      50ms     299ms
  235931 requests in 30.001s, net: in 35MB, out 31MB
HTTP Codes: 
     200       100.00%
Latency: 
                0s         6.76%
              10ms         3.20%
              20ms         5.05%
              30ms        13.20%
              40ms        23.83%
              50ms        21.30%
              60ms        11.15%
              70ms         5.78%
              80ms         3.27%
              90ms         2.23%
             100ms         1.49%
             110ms         1.01%
     120ms - 290ms         1.71%
Requests: 7864.10/sec
Net In: 9.2MBit/sec
Net Out: 8.2MBit/sec
Transfer: 2.2MB/sec
</pre>