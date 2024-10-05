document.addEventListener("DOMContentLoaded", () => {
  const list = document.querySelectorAll(".try-it-out");
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

  try {
    list.forEach((button) => {
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
  } catch (err) {
    console.error(err);
  }

  document.querySelectorAll(".return-from-try").forEach((button) => {
    button.addEventListener("click", () => {
      const tryIt = button.closest(".try-it");
      const routeInfo = tryIt.closest(".route-info");
      const exampleValue = routeInfo.querySelector(".example-value");
      const tryButton = routeInfo.querySelector(".try-it-out");
      const responseBlock = routeInfo.querySelector(".response");
      if (tryIt) {
        tryIt.style.display = "none";
        tryButton.style.display = "block";
        if (exampleValue) exampleValue.style.display = "block";
        if (responseBlock) {
          const responseExample =
            responseBlock.querySelector(".example-response");
          const responseResult =
            responseBlock.querySelector(".response-result");
          if (responseExample) responseExample.style.display = "block";
          if (responseResult) responseResult.style.display = "none";
        }
      }
    });
  });

  document.querySelectorAll(".type-selection select").forEach((select) => {
    select.value = select.options[0].value;

    select.addEventListener("change", () => {
      const handlerName = select.id.replace("content-type-", "");
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
      if (formDataInput) {
        formDataInput.style.display = "none";
        formDataEx.style.display = "none";
      }
    } else if (
      selectElement.value === "application/x-www-form-urlencoded" ||
      selectElement.value === "multipart/form-data"
    ) {
      if (jsonInput) {
        jsonInput.style.display = "none";
        jsonEx.style.display = "none";
      }
      formDataInput.style.display = "flex";
      formDataEx.style.display = "block";
    }
  }

  document.querySelectorAll(".request-form").forEach((form) => {
    form.addEventListener("submit", handleFormSubmit);
    const backButton = form.querySelector(".back-button");
    backButton.addEventListener("click", (event) => {
      event.preventDefault();
      const routeBlock = form.closest(".route-block");
      if (!routeBlock) {
        console.error("Route block not found");
        return;
      }

      const responseBlock = routeBlock.querySelector(".response");
      if (!responseBlock) {
        console.error("Response block not found");
        return;
      }

      const responseExample = responseBlock.querySelector(".example-response");
      const responseResult = responseBlock.querySelector(".response-result");

      if (responseExample) responseExample.style.display = "block";
      if (responseResult) responseResult.style.display = "none";

      backButton.style.display = "none";
    });
  });
});

async function handleFormSubmit(event) {
  event.preventDefault();
  const form = event.target;
  const routeBlock = form.closest(".route-block");
  if (!routeBlock) {
    console.error("Route block not found");
    return;
  }

  const backButton = form.querySelector(".back-button");
  if (!backButton) {
    console.error("Back button not found");
    return;
  }

  const routeHeader = routeBlock.querySelector(".route-header");
  if (!routeHeader) {
    console.error("Route header not found");
    return;
  }

  const method = routeHeader.querySelector(".method").textContent;
  if (!method) {
    console.error("Method not found");
    return;
  }

  const path = routeHeader.textContent.trim().split("\n")[1].trim();

  let responseBlock = routeBlock.querySelector(".response");
  if (!responseBlock) {
    responseBlock = document.createElement("div");
    responseBlock.classList.add("response");
    form.appendChild(responseBlock);
  }

  const responseExample = responseBlock.querySelector(".example-response");

  let responseDiv = responseBlock.querySelector(".response-result");
  if (!responseDiv) {
    responseDiv = document.createElement("div");
    responseDiv.classList.add("response-result");
    responseBlock.appendChild(responseDiv);
  }

  try {
    responseDiv.innerHTML = "<p>Loading...</p>";
    const response = await makeRequest(method, path, form, routeBlock);
    displayResponse(response, responseDiv);
    if (responseExample) responseExample.style.display = "none";
    backButton.style.display = "block";
  } catch (error) {
    displayError(error, responseDiv);
    if (responseExample) responseExample.style.display = "none";
    backButton.style.display = "block";
  }
}

async function makeRequest(method, path, form, routeBlock) {
  const baseUrl = window.location.origin;
  const fullPath = new URL(path, baseUrl);
  const acceptContentType = "application/json";
  const options = {
    method: method,
    headers: {
      Accept: acceptContentType,
    },
  };

  let typeSelection;
  let contentTypeSelect;
  if (routeBlock.querySelector(".type-selection")) {
    typeSelection = routeBlock.querySelector(".type-selection");
    contentTypeSelect = typeSelection.querySelector(
      'select[name="content-type"]'
    );
  }

  const contentType = contentTypeSelect ? contentTypeSelect.value : null;

  if (contentType === "application/json") {
    const jsonTextarea = form.querySelector('textarea[name="json-body"]');
    if (jsonTextarea && jsonTextarea.value) {
      options.headers["Content-Type"] = "application/json";
      try {
        options.body = JSON.stringify(JSON.parse(jsonTextarea.value));
      } catch (err) {
        throw new Error("Invalid JSON format");
      }
    } else {
      throw new Error("JSON body is required");
    }
  } else if (contentType === "multipart/form-data") {
    const formData = new FormData(form);
    options.body = formData;
  } else if (contentType === "application/x-www-form-urlencoded") {
    const formData = new FormData(form);
    options.headers["Content-Type"] = "application/x-www-form-urlencoded";
    options.body = new URLSearchParams(formData).toString();
  }

  const response = await fetch(fullPath, options);

  const contentTypeResponse = response.headers.get("content-type");
  if (contentTypeResponse && contentTypeResponse.includes("application/json")) {
    const jsonBody = await response.json();
    responseData = {
      body: jsonBody,
      headers: response.headers,
      url: response.url,
    };
  } else {
    const textBody = await response.text();
    responseData = {
      body: textBody,
      headers: response.headers,
      url: response.url,
    };
  }

  return responseData;
}

function displayResponse(responseData, responseDiv) {
  let responseMessage;
  if (typeof responseData.body === "object") {
    responseMessage = JSON.stringify(responseData.body, null, 2);
  } else {
    responseMessage = responseData.body;
  }

  let responseHTML = `
    <p style="margin-bottom:4px; ">Request URL:</p>
    <div class="request-url">
      <pre>${responseData.url}</pre>
    </div>
    <p style="margin-bottom:4px; ">Headers:</p>
    <div class="response-headers">
  `;
  for (const [key, value] of responseData.headers) {
    responseHTML += `
      <pre>${key}: ${value}</pre>
    `;
  }
  responseHTML += `
    </div>
    <p style="margin-bottom:4px; ">Response body:</p>
    <div class='code-block'>
       <pre>${responseMessage}</pre>
    </div>
  `;
  responseDiv.innerHTML = responseHTML;
  responseDiv.style.display = "block";
}

function displayError(error, responseDiv) {
  responseDiv.innerHTML = `
    <div class='code-block error'>
      <pre>${error.toString()}</pre>
    </div>
  `;
  responseDiv.style.display = "block";
}

function escapeHtml(unsafe) {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}
