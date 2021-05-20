Програма работает для одного уровня вложенности, любой из валидных ссылок http://target.true-tech.php.nixdev.co (Пока что нет валидации если стартовая ссылка нерабочая)

Что бы выполнть несколько уровней вложенности, нужно ввести основную ссылку на сайт (дальше сокращенные линки дописываютс к ней)

Для сайтов с небольшим количеством ссылок (к примеру http://vpustotu.ru), удается проверить все ссылки на нескольких уровнях

На сайте http://target.true-tech.php.nixdev.co для двух и более уровней вложенности полчучаю ошибку при отправлении запроса на сервер:
Get "http://target.true-tech.php.nixdev.co/valid/741": dial tcp 91.208.153.13:80: connectex: A connection attempt failed beсause the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond

Дл 3х b более уровней вложенности дополнительно возникает ошибка: Get "http://target.true-tech.php.nixdev.co/error/92674": dial tcp 91.208.153.13:80: connectex: Only one usage of each socket address (protocol/network address/port) is normally permitted. 

Мне кажется что я слишком активно пингую сервер и, возможно, интернет не позволяет все обработать, либо вместо метода гет существует какой то более легковесный способ достучатся до сервера. Пробовала добавлять таймауты, блокировать выполнение мьютексом, но не удалось исправить эту проблему
