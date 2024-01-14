Gstreamer-TestTask

Установка

git clone https://github.com/zanzhit/Gstreamer-TestTask
cd Gstreamer-TestTask
go mod download
Установите gstreamer.
Установите OpenVPN с нужным конфигом.

Использование

В файле urls.txt указать необходимые URL каждый с новой строки.
Запустите OpenVPN.
go run Gstreamer-TestTask.go в терминал.
Для прекращения записи нужно отправить что-нибудь в терминал (пустая строка, символ, предложение - неважно).
