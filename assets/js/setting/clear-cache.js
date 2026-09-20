'use strict';

$(document).ready(async function () {

    bindClearCacheEvents();
});

function bindClearCacheEvents() {

    let form = $("#clear-cache-view-form");

    form.off(".clearCache");

    form.on("click.clearCache", "#clear_cache_button", async function () {

        await saveClearCache();
    });
}

async function saveClearCache() {

    let confirmed = await confirmDialog("Clear All Cache", "Are you sure you want to clear all cached data?");

    if (!confirmed) {
        return;
    }

    let response = await web.post("/settings/clear-cache", new FormData());

    if (response.status === 200) {
        showAlertBox({title: "Success", message: "All cached data has been cleared successfully.", type: "success", sound: "/assets/sounds/success.mp3"});
    }
}

window.destroyPage = function () {

    $("#clear-cache-view-form").off(
        ".clearCache"
    );
};