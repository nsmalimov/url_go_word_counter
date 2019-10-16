#### Task

Программа читает из stdin строки, содержащие URL.
На каждый URL нужно отправить HTTP-запрос методом GET
и посчитать кол-во вхождений строки "Go" в теле ответа.
В конце работы приложение выводит на экран общее количество
найденных строк "Go" во всех переданных URL, например:

$ echo -e 'https://golang.org\nhttps://golang.org' | go run 1.go
Count for _https://golang.org_: 9
Count for https://golang.org: 9
Total: 18

Каждый URL должен начать обрабатываться сразу после вычитывания
и параллельно с вычитыванием следующего.
URL должны обрабатываться параллельно, но не более k=5 одновременно.
Обработчики URL не должны порождать лишних горутин, т.е. если k=5,
а обрабатываемых URL-ов всего 2, не должно создаваться 5 горутин.

Нужно обойтись без глобальных переменных и использовать только стандартную библиотеку.

Код необходимо залить в публичный репозиторий github или на gist.

#### Run tests

* Simple run: go test ./...

go test ./... -v (show prints)

* Show cover: go test -cover ./...

* Show full cover: go test ./... -coverprofile cover.out; go tool cover -func cover.out

* To html: go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

Покрытие: 95.8%

#### Run benchmarks


         