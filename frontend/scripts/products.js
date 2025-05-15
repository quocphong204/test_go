// Lấy token từ localStorage
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

  if (Array.isArray(data)) {
    data.forEach((p) => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>${p.name}</td>
        <td>${p.description}</td>
        <td>${p.price.toLocaleString()} đ</td>
      `;
      tbody.appendChild(row);
    });
  } else {
    tbody.innerHTML = `<tr><td colspan="3">Không có dữ liệu</td></tr>`;
  }
}

function logout() {
  localStorage.removeItem("token");
  window.location.href = "index.html";
}

fetchProducts();
