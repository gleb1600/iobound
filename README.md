# HTTP API, с помощью которого можно создавать продолжительный I/O bound задачи и получать результаты их работы.

## Работа

1. **Запуск**:
```bash
go run main.go
```
2. **Создание задачи, подача входных данные через поле input**
```bash
curl -X POST -d "{"input":"input_data_1"}" http://localhost:8080/create
```
*Ответ:*
```bash
Task started.
Task ID:a52ed698
```

3. **Просмотр задачи по ID**:
```bash
curl http://localhost:8080/a52ed698
```
**По умолчанию выполнение задачи занимает 30 секунд**

*Ответ, если задача не успела выполниться:*
```bash
Task ID:      a52ed698
Status:       in work
Input:        input_data_1
Output:       <no output yet>
Created at:   31.04.2025 12:24:48
```
*Ответ, если задача успела выполниться:*
```bash
Task ID:      a52ed698
Status:       completed
Input:        input_data_1
Output:       some result
Created at:   31.04.2025 12:24:48
```
4. **Просмотр всех задач**:
```bash
curl http://localhost:8080/tasks
```
*Ответ:*
```bash
Task ID:      21nj32yu
Status:       in work
Input:        input_data_3
Output:       <no output yet>
Created at:   31.04.2025 12:26:52

Task ID:      rg67d789
Status:       completed
Input:        input_data_2
Output:       some result
Created at:   31.04.2025 12:25:50

Task ID:      a52ed698
Status:       completed
Input:        input_data_1
Output:       some result
Created at:   31.04.2025 12:24:48
```
