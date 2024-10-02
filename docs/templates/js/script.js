document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll("section").forEach((section) => {
    section.querySelectorAll(".route-block").forEach((block) => {
      const header = block.querySelector(".route-header");
      const info = block.querySelector(".route-info");

      header.addEventListener("click", () => {
        if (info) {
          info.style.display =
            info.style.display === "none" || info.style.display === ""
              ? "block"
              : "none";
        }
      });
    });
  });

  document.querySelectorAll(".try-it-out").forEach((button) => {
    button.addEventListener("click", () => {
      const routeInfo = button.closest(".route-info");
      const exampleValue = routeInfo.querySelector(".example-value");
      const tryIt = routeInfo.querySelector(".try-it");

      if (tryIt) {
        if (exampleValue) {
          exampleValue.style.display = "none";
        }
        button.style.display = "none";
        tryIt.style.display = "block";
      }
    });
  });

  document.querySelectorAll(".return-from-try").forEach((button) => {
    button.addEventListener("click", () => {
      const tryIt = button.closest(".try-it");
      const routeInfo = tryIt.closest(".route-info");
      const exampleValue = routeInfo.querySelector(".example-value");
      const tryButton = routeInfo.querySelector(".try-it-out");
      if (tryIt) {
        tryIt.style.display = "none";
        tryButton.style.display = "block";
        exampleValue.style.display = "block";
      }
    });
  });

  function toggleInputFields(handlerName) {
    const selectElement = document.getElementById(
      `content-type-${handlerName}`
    );
    const jsonInput = document.getElementById(`json-input-${handlerName}`);
    const formDataInput = document.getElementById(
      `form-data-input-${handlerName}`
    );

    if (selectElement.value === "application/json") {
      jsonInput.style.display = "block";
      formDataInput.style.display = "none";
    } else if (selectElement.value === "application/x-www-form-urlencoded") {
      jsonInput.style.display = "none";
      formDataInput.style.display = "block";
    }
  }

  // Работа с формами
  const forms = document.querySelectorAll(".request-form");

  forms.forEach((form) => {
    const contentTypeSelect = form.querySelector(`select[name="content-type"]`);
    const jsonInput = form.querySelector(
      `#json-input-${form.id.split("-")[2]}`
    );
    const formDataInput = form.querySelector(
      `#form-data-input-${form.id.split("-")[2]}`
    );

    // Изначально показываем JSON и скрываем Form Data
    jsonInput.style.display = "block";
    formDataInput.style.display = "none";

    // Обработка смены типа контента
    contentTypeSelect.addEventListener("change", () => {
      if (contentTypeSelect.value === "application/json") {
        jsonInput.style.display = "block";
        formDataInput.style.display = "none";
      } else {
        jsonInput.style.display = "none";
        formDataInput.style.display = "block";
      }
    });

    // Обработка отправки формы
    form.addEventListener("submit", async (event) => {
      event.preventDefault(); // Предотвращение стандартной отправки формы

      let requestBody;
      const method = form
        .querySelector('button[type="submit"]')
        .classList.contains("post")
        ? "POST"
        : form
            .querySelector('button[type="submit"]')
            .classList.contains("delete")
        ? "DELETE"
        : form.querySelector('button[type="submit"]').classList.contains("get")
        ? "GET"
        : "GET";

      const action = form.action; // URL для запроса

      // Обработка JSON
      if (contentTypeSelect.value === "application/json") {
        requestBody = form.querySelector('textarea[name="json-body"]').value;
      } else {
        // Обработка Form Data
        const formData = new FormData(form);
        requestBody = formData; // Формат FormData
      }

      try {
        const response = await fetch(action, {
          method: method,
          headers:
            contentTypeSelect.value === "application/json"
              ? {
                  "Content-Type": "application/json",
                }
              : {},
          body:
            contentTypeSelect.value === "application/json"
              ? requestBody
              : requestBody, // Приведение JSON к строке
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const responseBody = await response.json();
        console.log("Response:", responseBody);
        alert("Response: " + JSON.stringify(responseBody, null, 2));
      } catch (error) {
        console.error("Error:", error);
        alert("Error: " + error.message);
      }
    });
  });
});
