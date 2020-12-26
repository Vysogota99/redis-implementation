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
<ul>
    <li>
        создание списка
        <code>
        curl -X POST
        -d '{"key":"asd"{"key":"list:1", "value":["ivan", 1, 3.2]}' http://127.0.0.1:3000/list/set
        </code>
    </li>
</ul>