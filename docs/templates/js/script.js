document.addEventListener("DOMContentLoaded", () => {
  console.log(
    "Количество .try-it-out элементов:",
    document.querySelectorAll(".try-it-out")
  );
  const list = document.querySelectorAll(".try-it-out");
  document.querySelectorAll("section").forEach((section) => {
    section.querySelectorAll(".route-block").forEach((block) => {
      const header = block.querySelector(".route-header");
      const info = block.querySelector(".route-info");
      console.log(block, "---------------");

      header.addEventListener("click", () => {
        console.log(block);

        if (info) {
          info.style.display =
            info.style.display === "none" || info.style.display === ""
              ? "block"
              : "none";
        }
      });
      // info.querySelectorAll(".type-selection select").forEach((select) => {
      //   const handlerName = select.id.replace("content-type-", "");
      //   console.log(handlerName);
      //   toggleInputFields(handlerName);
      // });
    });
  });
  console.log("Количество .try-it-out элементов:", list);

  try {
    list.forEach((button) => {
      console.log(button, "sadsa");
      button.addEventListener("click", () => {
        console.log(button, "sadsa");
        const routeInfo = button.closest(".route-info");
        const exampleValue = routeInfo.querySelector(".example-value");
        const tryIt = routeInfo.querySelector(".try-it");

        console.log(tryIt);

        if (tryIt) {
          if (exampleValue) {
            exampleValue.style.display = "none";
          }
          button.style.display = "none";
          tryIt.style.display = "block";
        }
      });
    });
  } catch (err) {
    console.error(err);
  }

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

  document.querySelectorAll(".type-selection select").forEach((select) => {
    select.value = select.options[0].value;

    select.addEventListener("change", () => {
      const handlerName = select.id.replace("content-type-", "");
      console.log(handlerName);
      toggleInputFields(handlerName);
    });

    const event = new Event("change");
    select.dispatchEvent(event);
  });

  function toggleInputFields(handlerName) {
    const selectElement = document.getElementById(
      `content-type-${handlerName}`
    );
    const jsonInput = document.getElementById(`json-input-${handlerName}`);
    const formDataInput = document.getElementById(
      `form-data-input-${handlerName}`
    );
    const jsonEx = document.querySelector(`.json-example-${handlerName}`);
    const formDataEx = document.querySelector(
      `.form-data-example-${handlerName}`
    );
    if (!selectElement) {
      console.error(`Element #content-type-${handlerName} not found`);
      return;
    }

    if (selectElement.value === "application/json") {
      jsonInput.style.display = "block";
      jsonEx.style.display = "block";
      formDataInput.style.display = "none";
      formDataEx.style.display = "none";
    } else if (
      selectElement.value === "application/x-www-form-urlencoded" ||
      selectElement.value === "multipart/form-data"
    ) {
      jsonInput.style.display = "none";
      jsonEx.style.display = "none";
      formDataInput.style.display = "flex";
      formDataEx.style.display = "block";
    }
  }

  // window.toggleInputFields = toggleInputFields;
});
