# ModbusSlave

Включается TCP Server по 502 порту и ждет подключения.
При включении Modbus TCP Master (TCP Client) посылается телеграмма на TCP Server и дальше идет обработка. Слайс regSlice - слайс из значений.
Выбор SlaveID в глобальной переменной, только для функции 03 Holding Registor, любое Quantity, но меньше длины regSlice.
