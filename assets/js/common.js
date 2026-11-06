if (typeof bootstrap !== 'undefined' && bootstrap.Modal && bootstrap.Modal.Default) {
    bootstrap.Modal.Default.backdrop = 'static';
    bootstrap.Modal.Default.keyboard = false;
}

const MODAL_BASE_ZINDEX = 1050;
const MODAL_STEP_ZINDEX = 20;

function reorderModalStack() {
    let $openModals = $('.modal.show');
    if ($openModals.length === 0) return;

    $('body').addClass('modal-open');
    document.body.style.overflow = 'hidden';

    let $backdrops = $('.modal-backdrop');
    $openModals.each(function (index) {
        let baseZ = MODAL_BASE_ZINDEX + (index * MODAL_STEP_ZINDEX);
        $(this).css('z-index', baseZ + 5);
        if ($backdrops.eq(index).length) {
            $backdrops.eq(index).css('z-index', baseZ);
        }
    });
}

function getAvailablePopupContainer(preferredDiv) {
    let activeContainers = new Set();
    document.querySelectorAll('.modal.show').forEach(function (m) {
        let parent = m.closest('[id^="popup_"]');
        if (parent && parent.id) {
            activeContainers.add(parent.id);
        }
    });

    if (preferredDiv && preferredDiv !== 'content') {
        if (!activeContainers.has(preferredDiv)) {
            let el = document.getElementById(preferredDiv);
            if (!el || !el.querySelector('.modal.show')) {
                return preferredDiv;
            }
        }
    }

    let i = 1;
    while (activeContainers.has('popup_' + i)) {
        i++;
    }
    return 'popup_' + i;
}

$(document).on('show.bs.modal', '.modal', function () {
    $('body').addClass('modal-open');
    let openCount = $('.modal.show').not(this).length;
    let baseZ = MODAL_BASE_ZINDEX + (openCount * MODAL_STEP_ZINDEX);
    $(this).css('z-index', baseZ + 5);

    setTimeout(function () {
        reorderModalStack();
    }, 0);
});

$(document).on('shown.bs.modal', '.modal', function () {
    reorderModalStack();
});

$(document).on('hidden.bs.modal', '.modal', function () {
    let $thisModal = $(this);
    setTimeout(function () {
        let $openModals = $('.modal.show');
        if ($openModals.length > 0) {
            $('body').addClass('modal-open');
            document.body.style.overflow = 'hidden';
            reorderModalStack();
            $openModals.last().focus();
        } else {
            $('body').removeClass('modal-open');
            document.body.style.overflow = '';
        }

        let $container = $thisModal.closest('[id^="popup_"]');
        if ($container.length && $container.attr('id') !== 'popup_1') {
            if ($container.find('.modal.show').length === 0) {
                let containerId = $container.attr('id');
                if (typeof containerAssets !== 'undefined' && containerAssets.has(containerId)) {
                    containerAssets.delete(containerId);
                }
                $container.remove();
            }
        }
    }, 10);
});

function showAlertBox(options = {}) {
    let {title = "Loading", message = "Please wait...", type = "theme", delay = 2000, timer = 300, positionX = "right", positionY = "top", sound = null, animation = "viho-up", style = "bootstrap",} = options;

    if (sound) {
        let audio = new Audio(sound);
        audio.play().catch(() => {});
    }

    let html = title + ' | ' + message;

    return $.notify(html, {type: type, allow_dismiss: true, delay: delay, timer: timer, showProgressbar: true, style: style, placement: {from: positionY, align: positionX}, z_index: 999999, animate: {enter: "animated " + animation, exit: "animated fadeOut"}});
}

function openPopupContainer(cleanedHtml, title, popupWidth, contentDiv, targetModalSelector) {
    contentDiv = getAvailablePopupContainer(contentDiv);
    let $container = $('#' + contentDiv);
    if ($container.length === 0) {
        $('body').append(`<div id="${contentDiv}" class="popup-stack-container"></div>`);
        $container = $('#' + contentDiv);
    }

    let temp = document.createElement('div');
    temp.innerHTML = cleanedHtml;
    let modalEl = temp.querySelector('.modal');

    if (modalEl) {
        $container.html(cleanedHtml);

        let targetSelector = targetModalSelector;
        if (!targetSelector && title && title.toLowerCase().includes('activity')) {
            targetSelector = '#modal_designation_activity, .modal[id*="activity"]';
        }
        let liveModal = targetSelector ? document.querySelector('#' + contentDiv + ' ' + targetSelector) : null;
        if (!liveModal) {
            liveModal = document.querySelector('#' + contentDiv + ' .modal');
        }
        if (liveModal) {
            let bsModal = bootstrap.Modal.getOrCreateInstance(liveModal);
            bsModal.show();
        }
    } else {
        let modalId = contentDiv + '_modal';
        let $existingModal = $('#' + modalId);
        let sizeClass = popupWidth >= 1140 ? 'modal-xl' : (popupWidth >= 768 ? 'modal-lg' : '');

        if ($existingModal.length === 0) {
            let modalHtml = `
                <div class="modal fade" id="${modalId}" tabindex="-1" aria-hidden="true" data-bs-backdrop="static">
                    <div class="modal-dialog modal-dialog-centered ${sizeClass}">
                        <div class="modal-content rounded-3 border-0 shadow-lg overflow-hidden">
                            <div class="modal-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
                                <h5 class="modal-title fw-bold text-dark fs-15 mb-0" id="${modalId}_title">${title || ''}</h5>
                                <button type="button" class="btn-close custom-btn-close shadow-none" data-bs-dismiss="modal" aria-label="Close"><i class="ti ti-x"></i></button>
                            </div>
                            <div class="modal-body p-4 popup-content-body">
                            </div>
                        </div>
                    </div>
                </div>`;
            $container.html(modalHtml);
            $existingModal = $('#' + modalId);
        } else {
            $('#' + modalId + '_title').text(title || '');
        }

        $existingModal.find('.popup-content-body').html(cleanedHtml);
        let bsModal = bootstrap.Modal.getOrCreateInstance($existingModal[0]);
        bsModal.show();
    }
}

$(document).on("click", ".menu-link", function (e) {
    e.preventDefault();
    let $this = $(this);
    if ($this.data("loading")) return;
    $this.data("loading", true);
    let url = $this.data("url");
    let title = $this.data("title");
    let data = $this.data("data") || "";
    let name = $this.data("name") || $this.data("content-action") || $this.data("action") || "";
    let popupStatus = Number($this.data("popup-status"));
    let popupWidth = Number($this.data("popup-width"));
    let contentDiv = $this.data("content-div");

    openMenu(url, title, data, name, popupStatus, popupWidth, contentDiv).finally(function () {
        $this.removeData("loading");
    });
});


async function openMenu(url, title, data, name, popupStatus, popupWidth, contentDiv) {

    let systemName = document.querySelector('meta[name="system_name"]')?.getAttribute('content');

    if (title && systemName) {
        document.title = title + ' | ' + systemName;
    }

    if (popupStatus === 1) {
        contentDiv = getAvailablePopupContainer(contentDiv);
    }

    let response = await web.get(data ? `${url}/${data}` : url);

    let result = response.data;

    let {cleanedHtml, cssUrls, jsUrls} = extractAndLoadAssets(result.data);

    if (popupStatus === 1) {
        openPopupContainer(cleanedHtml, title, popupWidth, contentDiv);
    }
    else {
        unloadContainerAssets(contentDiv);
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