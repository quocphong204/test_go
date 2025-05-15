const token = localStorage.getItem("token");
if (!token) {
  alert("Bạn chưa đăng nhập!");
  window.location.href = "index.html";
}

let role = "user";

async function fetchOrders() {
  const res = await fetch("http://localhost:8080/api/orders", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  const data = await res.json();
  const tbody = document.querySelector("#ordersTable tbody");
  tbody.innerHTML = "";

  if (Array.isArray(data)) {
    data.forEach((order) => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>#${order.id}</td>
        <td>${order.total_price.toLocaleString()} đ</td>
        <td>${order.status}</td>
        <td>${order.payment_status}</td>
        <td id="actions_${order.id}"></td>
      `;
      tbody.appendChild(row);

      const actionsTd = document.getElementById(`actions_${order.id}`);

      // Nếu chưa thanh toán
      if (order.payment_status === "unpaid") {
        const payBtn = document.createElement("button");
        payBtn.innerText = "Thanh toán";
        payBtn.onclick = () => markPaid(order.id);
        actionsTd.appendChild(payBtn);
      }

      // Nếu là admin → cho xử lý đơn
      if (role === "admin" && order.status === "pending") {
        const processBtn = document.createElement("button");
        processBtn.innerText = "Xử lý";
        processBtn.onclick = () => processOrder(order.id);
        actionsTd.appendChild(processBtn);
      }
    });
  }
}

function logout() {
  localStorage.removeItem("token");
  window.location.href = "index.html";
}

async function fetchUserInfo() {
  const res = await fetch("http://localhost:8080/api/profile", {
    headers: { Authorization: `Bearer ${token}` },
  });
  const data = await res.json();
  role = data.role || "user";
  fetchOrders();
}

async function markPaid(orderId) {
  if (!confirm("Xác nhận thanh toán cho đơn #" + orderId + "?")) return;

  const res = await fetch(`http://localhost:8080/api/orders/${orderId}/pay`, {
    method: "PUT",
    headers: { Authorization: `Bearer ${token}` },
  });

  const data = await res.json();
  alert(data.message || "✅ Đã thanh toán");
  fetchOrders();
}

async function processOrder(orderId) {
  if (!confirm("Xác nhận xử lý đơn #" + orderId + "?")) return;

  const res = await fetch(`http://localhost:8080/api/orders/${orderId}/process`, {
    method: "PUT",
    headers: { Authorization: `Bearer ${token}` },
  });

  const data = await res.json();
  alert(data.message || "✅ Đã xử lý đơn");
  fetchOrders();
}

fetchUserInfo();
