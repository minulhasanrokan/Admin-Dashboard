'use strict';

$(document).ready(function () {
    bindSmsSettingsEvents();
});

function bindSmsSettingsEvents() {
    let form = $("#sms-settings-edit-form");

    form.off(".smsSettings");

    $("#toggle_token_visibility").off("click").on("click", function () {
        let input = $("#ssl_api_token");
        let icon = $("#token_eye_icon");
        if (input.attr("type") === "password") {
            input.attr("type", "text");
            icon.removeClass("ti-eye").addClass("ti-eye-off");
        } else {
            input.attr("type", "password");
            icon.removeClass("ti-eye-off").addClass("ti-eye");
        }
    });

    form.on("focusout.smsSettings", "#ssl_api_url", function () {
        form_validation({ control: "ssl_api_url#text#mandatory", msg_text: "API URL" });
    });

    form.on("focusout.smsSettings", "#ssl_api_token", function () {
        form_validation({ control: "ssl_api_token#text#mandatory", msg_text: "API Token" });
    });

    form.on("focusout.smsSettings", "#ssl_sid", function () {
        form_validation({ control: "ssl_sid#text#mandatory", msg_text: "SID" });
    });

    form.on("click.smsSettings", "#sms_setting_button", async function () {
        await saveSmsSettings();
    });

    $("#send_test_sms_btn").off("click").on("click", async function () {
        await sendTestSms();
    });
}

async function saveSmsSettings() {
    let status = form_validation({
        control: "ssl_api_url#text#mandatory*ssl_api_token#text#mandatory*ssl_sid#text#mandatory",
        msg_text: "API URL*API Token*SID"
    });

    if (status === false) {
        showAlertBox({ title: "Error", message: "Validation failed. Please fill all required fields.", type: "danger", sound: "/assets/sounds/error.mp3" });
        return false;
    }

    let confirmed = await confirmDialog("SMS Settings", "Are you sure you want to save the SMS API Settings?");
    if (!confirmed) {
        return;
    }

    showLoader();

    let formData = new FormData();
    formData.append("sms_driver", $("#sms_driver").val() || "sslcommerz");
    formData.append("status", $("#status").is(":checked") ? "1" : "0");
    formData.append("ssl_api_url", $("#ssl_api_url").val());
    formData.append("ssl_api_token", $("#ssl_api_token").val());
    formData.append("ssl_sid", $("#ssl_sid").val());
    formData.append("ssl_csms_id_prefix", $("#ssl_csms_id_prefix").val());

    let response = await web.post('/settings/sms-settings', formData);

    if (response.status === 200) {
        showAlertBox({
            title: "Success",
            message: "SMS API Settings updated successfully",
            type: "info",
            sound: "/assets/sounds/success.mp3"
        });
    }
}

async function sendTestSms() {
    let mobile = $("#test_mobile").val().trim();
    let message = $("#test_message").val().trim();
    let alertBox = $("#test_sms_alert");

    alertBox.addClass("d-none").removeClass("alert-success alert-danger").text("");

    if (!mobile) {
        $("#test_mobile").addClass("is-invalid");
        $("#error_msg_test_mobile").text("Recipient Mobile Number is required").show();
        return;
    } else {
        $("#test_mobile").removeClass("is-invalid");
        $("#error_msg_test_mobile").text("").hide();
    }

    if (!message) {
        $("#test_message").addClass("is-invalid");
        $("#error_msg_test_message").text("Message content is required").show();
        return;
    } else {
        $("#test_message").removeClass("is-invalid");
        $("#error_msg_test_message").text("").hide();
    }

    let btn = $("#send_test_sms_btn");
    let originalHtml = btn.html();
    btn.prop("disabled", true).html('<span class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span> Sending...');

    let formData = new FormData();
    formData.append("mobile", mobile);
    formData.append("message", message);

    try {
        let response = await web.post('/settings/sms-settings/test', formData);

        if (response.status === 200) {
            alertBox.removeClass("d-none").addClass("alert-success").text(response.data.message || "Test SMS sent successfully!");
            showAlertBox({
                title: "Test SMS Sent",
                message: "Test SMS dispatched to " + mobile,
                type: "info",
                sound: "/assets/sounds/success.mp3"
            });
        } else {
            let errorMsg = response.data && response.data.message ? response.data.message : "Failed to send test SMS";
            alertBox.removeClass("d-none").addClass("alert-danger").text(errorMsg);
        }
    } catch (err) {
        let errorMsg = "An unexpected error occurred while sending test SMS";
        if (err.response && err.response.data && err.response.data.message) {
            errorMsg = err.response.data.message;
        }
        alertBox.removeClass("d-none").addClass("alert-danger").text(errorMsg);
    } finally {
        btn.prop("disabled", false).html(originalHtml);
    }
}

window.destroyPage = function () {
    $("#sms-settings-edit-form").off(".smsSettings");
    $("#send_test_sms_btn").off("click");
    $("#toggle_token_visibility").off("click");
};
