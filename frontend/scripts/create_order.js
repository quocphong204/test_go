const token = localStorage.getItem("token");
if (!token) {
  alert("Bạn chưa đăng nhập!");
  window.location.href = "index.html";
}

async function fetchProducts() {
  const res = await fetch("http://localhost:8080/api/products", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  const data = await res.json();

  const tbody = document.querySelector("#productTable tbody");
  tbody.innerHTML = "";

  data.forEach((p) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${p.name}</td>
      <td>${p.price.toLocaleString()} đ</td>
      <td>
        <input type="number" min="0" value="0" name="quantity_${p.id}" data-product-id="${p.id}" />
      </td>
    `;
    tbody.appendChild(row);
  });
}

document.getElementById("orderForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const inputs = document.querySelectorAll("input[type='number']");
  const items = [];

  inputs.forEach((input) => {
    const qty = parseInt(input.value);
    if (qty > 0) {
      items.push({
        product_id: parseInt(input.dataset.productId),
        quantity: qty,
      });
    }
  });

  if (items.length === 0) {
    alert("Vui lòng chọn ít nhất 1 sản phẩm!");
    return;
  }

  const res = await fetch("http://localhost:8080/api/orders", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ items }),
  });

  const data = await res.json();
  if (res.ok) {
    alert("✅ Đơn hàng đã được tạo!");
    window.location.href = "orders.html";
  } else {
    alert("❌ Lỗi tạo đơn: " + data.error);
  }
});

fetchProducts();
