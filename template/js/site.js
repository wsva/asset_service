function init_button_login() {
  $("#oauth2-login").href = `/oauth2/login?return_to=${encodeURIComponent(
    window.location.href
  )}`;
}

function init_button_password() {
  $("#btn-encrypt").href = `/asset/password?action=encrypt`;
  $("#btn-decrypt").href = `/asset/password?action=decrypt`;
}

function init_password_overview() {
  $.getJSON("/asset/password?action=overview", function (data) {
    if (!data.success) {
      return;
    }
    data.data.list.forEach(function (item) {
      switch (item.name) {
        case "encrypted":
          $("#encrypted-count").text(item.value);
          break;
        case "unencrypted":
          $("#unencrypted-count").text(item.value);
          break;
      }
    });
  });
}

function init_password_project() {
  $.getJSON(
    `/asset/password?action=project&account_id=${encodeURIComponent(
      $("#account-email").text()
    )}`,
    function (data) {
      if (!data.success) {
        return;
      }
      const $select = $("#config-project");
      $select.empty();
      $select.append(`<option value="">select project first</option>`);
      data.data.list.forEach(function (item) {
        $select.append(`<option value="${item[0]}">${item[1]}</option>`);
      });
    }
  );
}

function init_password_description(project_id) {
  $.getJSON(
    `/asset/password?action=description&account_id=${encodeURIComponent(
      $("#account-email").text()
    )}&project_id=${encodeURIComponent(project_id)}`,
    function (data) {
      console.log(data);
      if (!data.success) {
        return;
      }
      console.log(data);
      const $result = $("#config-result");
      $result.empty();
      if (!Array.isArray(data.data.list)) {
        $result.append(`<div>no data found</div>`);
        return;
      }
      data.data.list.forEach(function (item) {
        $result.append(`<div class="result-item" uuid="${item[0]}">
        <div class="result-item-content">${item[1]}</div>
        </div>`);
      });
    }
  );
}

function init_ip() {
  $.getJSON(
    `/asset/ip`,
    function (data) {
      console.log(data);
      if (!data.success) {
        return;
      }
      console.log(data);
      const $result = $("#config-result");
      $result.empty();
      if (!Array.isArray(data.data.list)) {
        $result.append(`<div>no data found</div>`);
        return;
      }
      data.data.list.forEach(function (item) {
        $result.append(`<div class="result-item">
        <div class="result-item-content">${item[0]} ${item[1]}</div>
        </div>`);
      });
    }
  );
}

function init_address() {
  $.getJSON(
    `/asset/address`,
    function (data) {
      console.log(data);
      if (!data.success) {
        return;
      }
      console.log(data);
      const $result = $("#config-result");
      $result.empty();
      if (!Array.isArray(data.data.list)) {
        $result.append(`<div>no data found</div>`);
        return;
      }
      data.data.list.forEach(function (item) {
        $result.append(`<div class="result-item">
        <div class="result-item-content">${item[0]} ${item[1]} ${item[2]}</div>
        </div>`);
      });
    }
  );
}

function init_tab_password() {
    $("#config-result").empty();
    init_password_overview();
    init_password_project();
}

function init_tab_ip() {
    $("#config-result").empty();
    init_ip();
}

function init_tab_address() {
    $("#config-result").empty();
    init_address();
}