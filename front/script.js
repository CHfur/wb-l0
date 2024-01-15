function getOrderDetails() {
    // Получаем значение Order ID из формы
    var orderId = document.getElementById('orderId').value;

    // Проверяем, что orderId не пуст
    if (orderId.trim() === "") {
        alert("Please enter an Order ID");
        return;
    }

    // Формируем URL для запроса данных о заказе
    var apiUrl = 'http://localhost:8080/v1/order/' + orderId;

    // Отправляем GET-запрос на сервер
    fetch(apiUrl)
        .then(response => response.json())
        .then(data => displayOrderDetails(data))
        .catch(error => console.error('Error:', error));
}

function displayOrderDetails(orderDetails) {
    // Отображаем полученные данные о заказе
    var orderDetailsDiv = document.getElementById('orderDetails');
    orderDetailsDiv.innerHTML = '<pre>' + JSON.stringify(orderDetails, null, 2) + '</pre>';
}