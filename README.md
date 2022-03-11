# ModbusSlave

Включается TCP Server по 502 порту и ждет подключения.
При включении Modbus TCP Master (TCP Slave) посылается телеграмма на TCP Server и дальше идет обработка. Слайс regSlice - слайс из значений.
Работает для любого Slave ID, только для функции 03 Holding Registor, любое Quantity, но меньше длины regSlice.
