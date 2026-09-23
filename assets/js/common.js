if (typeof bootstrap !== 'undefined' && bootstrap.Modal && bootstrap.Modal.Default) {
    bootstrap.Modal.Default.backdrop = 'static';
    bootstrap.Modal.Default.keyboard = false;
}

function showAlertBox(options = {}) {
    let {title = "Loading", message = "Please wait...", type = "theme", delay = 2000, timer = 300, positionX = "right", positionY = "top", sound = null, animation = "viho-up", style = "bootstrap",} = options;

    if (sound) {
        let audio = new Audio(sound);
        audio.play().catch(() => {});
    }

    let html = title + ' | ' + message;

    return $.notify(html, {type: type, allow_dismiss: true, delay: delay, timer: timer, showProgressbar: true, style: style, placement: {from: positionY, align: positionX}, z_index: 999999, animate: {enter: "animated " + animation, exit: "animated fadeOut"}});
}

$(document).on("click", ".menu-link", function (e) {
    e.preventDefault();
    let $this = $(this);
    openMenu($this.data("url"), $this.data("title"), $this.data("data") || "", $this.data("name"), Number($this.data("popup-status")), Number($this.data("popup-width")), $this.data("content-div")).then();
});

$(document).on("click", ".btn-grid-menu-action", function (e) {
    e.preventDefault();
    let $this = $(this);
    let url = $this.data("url");
    let title = $this.data("title");
    let popupStatus = Number($this.data("popup-status"));
    let width = Number($this.data("popup-width"));
    if (typeof window.get_new_page === "function") {
        window.get_new_page(url, title, "", "", popupStatus, width);
    }
});

async function openMenu(url, title, data, name, popupStatus, popupWidth, contentDiv) {

    let systemName = document.querySelector('meta[name="system_name"]')?.getAttribute('content');

    document.title = title + ' | ' + systemName;

    let response = await web.get(data ? `${url}/${data}` : url);

    let result = response.data;

    let {cleanedHtml, cssUrls, jsUrls} = extractAndLoadAssets(result.data);

    unloadContainerAssets(contentDiv);

    if (popupStatus === 1) {
    }
    else {

        $('#' + contentDiv).html(cleanedHtml);
    }

    await Promise.all(cssUrls.map(function (href) {return loadCSS(href);}));

    for (let src of jsUrls) {

        await loadScript(src);
    }

    containerAssets.set(contentDiv, {css: cssUrls, js: jsUrls});
}

async function loadDropDown(tableName, fieldName, id, container, title, dropDownType, selectedValue = '', disableStatus = 0, tabIndex = '', otherField = "", otherValue = "", notEqualField = 0, notEqualValue = 0) {

    let value = selectedValue || 0;
    otherField = otherField || 0;
    otherValue = otherValue || 0;
    notEqualField = notEqualField || 0;
    notEqualValue = notEqualValue || 0;

    let response = await web.get("common/get-dropdown", {
        params: {
            table_name: tableName,
            field_name: fieldName,
            value: value,
            other_field: otherField,
            other_value: otherValue,
            not_equal_field: notEqualField,
            not_equal_value: notEqualValue
        }
    });

    let result = response.data;

    let html = "";

    if (dropDownType === 1) {
        html += `<select tabindex="${tabIndex}" name="${id}" id="${id}" class="form-control select2" multiple="multiple">`;
    }
    else {
        html += `<select tabindex="${tabIndex}" name="${id}" id="${id}" class="form-control select2">`;
    }

    html += `<option value="">${title}</option>`;
    let fieldNameArr = fieldName.split(",");

    if (result && result.success && Array.isArray(result.data)) {
        result.data.forEach(item => {
            html += `<option value="${item[fieldNameArr[0].trim()]}">${item[fieldNameArr[1].trim()]}</option>`;
        });
    }

    html += `</select>`;

    document.getElementById(container).innerHTML = html;

    if (fieldNameArr.length > 2 && result && Array.isArray(result.data)) {
        let hiddenFields = "";
        result.data.forEach(item => {
            hiddenFields += `<input type="hidden" name="hidden_field_${fieldNameArr[2].trim()}-${item[fieldNameArr[0].trim()]}" id="hidden_field_${fieldNameArr[2].trim()}-${item[fieldNameArr[0].trim()]}" value="${item[fieldNameArr[2].trim()]}">`;
        });
        document.getElementById(container).insertAdjacentHTML("beforeend", hiddenFields);
    }

    let $select = $("#" + id);

    if ($select.hasClass("select2-hidden-accessible")) {
        $select.select2("destroy");
    }

    $select.select2({
        //placeholder: title,
        //allowClear: true,
        width: "100%"
    });

    if (selectedValue !== "") {

        if (dropDownType === 1) {
            $select.val(selectedValue.split(","));
        }
        else {
            $select.val(selectedValue);
        }
    }
    else {
        $select.val(null);
    }

    $select.trigger("change");

    if (disableStatus === 1) {
        $select.prop("disabled", true);
    }
}

function blankDropDown(id, container, title, dropDownType, tabIndex = '', disableStatus = 0) {

    let html = "";

    if (dropDownType === 1) {
        html += `<select tabindex="${tabIndex}" name="${id}" id="${id}" class="form-control select2" multiple="multiple">`;
    }
    else {
        html += `<select tabindex="${tabIndex}" name="${id}" id="${id}" class="form-control select2">`;
    }

    html += `<option value="">${title}</option>`;
    html += `</select>`;

    document.getElementById(container).innerHTML = html;

    let $select = $("#" + id);

    if ($select.hasClass("select2-hidden-accessible")) {
        $select.select2("destroy");
    }

    $select.select2({
        placeholder: title,
        allowClear: true,
        width: "100%"
    });

    $select.val(null).trigger("change");

    if (disableStatus === 1) {
        $select.prop("disabled", true).trigger("change");
    }
}

$(document).on('focus', '.select2-selection', function () {
    $(this).closest('.select2-container').prev('select').select2('open');
});

$(document).on("click", ".logout-link", function (e) {
    e.preventDefault();
    let url = $(this).attr("href");
    if (typeof confirmDialog === "function") {
        confirmDialog("Logout", "Are you sure you want to log out?", "warning").then(function (confirmed) {
            if (confirmed) {
                window.location.href = url;
            }
        });
    } else {
        window.location.href = url;
    }
});