'use strict';

$(document).ready(function () {
    initEmailSelect2();
    bindEmailSettingsEvents();
});

function initEmailSelect2() {
    if (!$.fn.select2) {
        return;
    }

    let $smtpEnc = $("#smtp_encryption");
    if ($smtpEnc.length > 0) {
        if ($smtpEnc.hasClass("select2-hidden-accessible")) {
            $smtpEnc.select2("destroy");
        }
        $smtpEnc.select2({
            dropdownParent: $("#modal_smtp"),
            //placeholder: "Select Encryption",
            //allowClear: true,
            width: "100%"
        });
    }

    let $awsReg = $("#aws_region");
    if ($awsReg.length > 0) {
        if ($awsReg.hasClass("select2-hidden-accessible")) {
            $awsReg.select2("destroy");
        }
        $awsReg.select2({
            dropdownParent: $("#modal_ses"),
            //placeholder: "Select AWS Region",
            //allowClear: true,
            width: "100%"
        });
    }

    $(document).off("click.labelAwsRegion", "label[for='aws_region']").on("click.labelAwsRegion", "label[for='aws_region']", function (e) {
        e.preventDefault();
        $("#aws_region").select2("open");
    });

    $(document).off("click.labelSmtpEnc", "label[for='smtp_encryption']").on("click.labelSmtpEnc", "label[for='smtp_encryption']", function (e) {
        e.preventDefault();
        $("#smtp_encryption").select2("open");
    });
}

function bindEmailSettingsEvents() {
    let form = $("#email-settings-edit-form");

    form.off(".emailSettings");

    form.on("focusout.emailSettings", "#from_email", function () {
        form_validation({ control: "from_email#email#mandatory", msg_text: "From Email" });
    });

    form.on("focusout.emailSettings", "#from_name", function () {
        form_validation({ control: "from_name#text#mandatory", msg_text: "From Name" });
    });

    $('input[name="mail_driver"]').off("change.emailDriver").on("change.emailDriver", function () {
        updateActiveCardHighlight();
    });

    form.on("click.emailSettings", "#email_setting_button", async function () {
        await saveEmailSettings();
    });

    $(".toggle-password").off("click").on("click", function () {
        let targetId = $(this).data("target");
        let input = $(targetId);
        let icon = $(this).find("i");
        if (input.attr("type") === "password") {
            input.attr("type", "text");
            icon.removeClass("ti-eye").addClass("ti-eye-off");
        } else {
            input.attr("type", "password");
            icon.removeClass("ti-eye-off").addClass("ti-eye");
        }
    });

    $("#save_sendgrid_btn").off("click").on("click", async function () {
        await saveSendGridSettings();
    });

    $("#save_smtp_btn").off("click").on("click", async function () {
        await saveSmtpSettings();
    });

    $("#save_ses_btn").off("click").on("click", async function () {
        await saveSesSettings();
    });

    $("#aws_region").off("change.emailRegion").on("change.emailRegion", function () {
        let val = $(this).val();
        if (val) {
            $(this).removeClass("is-invalid");
            $(this).next(".select2-container").find(".select2-selection").removeClass("custom-is-invalid is-invalid");
            $("#error_msg_aws_region").text("").hide();
        }
    });

    $("#smtp_encryption").off("change.emailSmtpEnc").on("change.emailSmtpEnc", function () {
        let val = $(this).val();
        if (val) {
            $(this).removeClass("is-invalid");
            $(this).next(".select2-container").find(".select2-selection").removeClass("custom-is-invalid is-invalid");
            $("#error_msg_smtp_encryption").text("").hide();
        }
    });

    $("#send_test_email_btn").off("click").on("click", async function () {
        await sendTestEmail();
    });
}

function updateActiveCardHighlight() {
    let selectedDriver = $('input[name="mail_driver"]:checked').val();
    $("#card-sendgrid, #card-smtp, #card-ses").removeClass("border-primary");
    if (selectedDriver) {
        $("#card-" + selectedDriver).addClass("border-primary");
    }
}

async function saveEmailSettings() {
    let status = form_validation({
        control: "from_email#email#mandatory*from_name#text#mandatory",
        msg_text: "From Email*From Name"
    });

    let mailDriver = $('input[name="mail_driver"]:checked').val();
    if (!mailDriver) {
        $("#error_msg_mail_driver").text("Please select a Mail Driver").show();
        status = false;
    } else {
        $("#error_msg_mail_driver").text("").hide();
    }

    if (status === false) {
        showAlertBox({
            title: "Error",
            message: "Validation failed. Please fill all required fields.",
            type: "danger",
            sound: "/assets/sounds/error.mp3"
        });
        return false;
    }

    let confirmed = await confirmDialog("Email Settings", "Are you sure you want to save the Email Settings?");
    if (!confirmed) {
        return;
    }

    let formData = new FormData();
    formData.append("from_email", $("#from_email").val().trim());
    formData.append("from_name", $("#from_name").val().trim());
    formData.append("mail_driver", mailDriver);

    let response = await web.post('/settings/email-settings', formData);

    if (response.status === 200) {
        showAlertBox({
            title: "Success",
            message: "Email settings updated successfully",
            type: "info",
            sound: "/assets/sounds/success.mp3"
        });
        updateActiveCardHighlight();
    }
}

async function saveSendGridSettings() {
    let apiKey = $("#sendgrid_api_key").val().trim();

    if (!apiKey) {
        $("#sendgrid_api_key").addClass("is-invalid");
        $("#error_msg_sendgrid_api_key").text("SendGrid API Key is required").show();
        return;
    } else {
        $("#sendgrid_api_key").removeClass("is-invalid");
        $("#error_msg_sendgrid_api_key").text("").hide();
    }

    let btn = $("#save_sendgrid_btn");
    let originalHtml = btn.html();
    btn.prop("disabled", true).html('<span class="spinner-border spinner-border-sm me-1" role="status"></span> Saving...');

    try {
        let formData = new FormData();
        formData.append("sendgrid_api_key", apiKey);

        let response = await web.post('/settings/email-settings/sendgrid', formData);

        if (response.status === 200) {
            $("#modal_sendgrid").modal("hide");
            $("#sendgrid-status-badge")
                .removeClass("bg-light text-muted border fw-normal")
                .addClass("bg-success-transparent text-success fw-semibold")
                .html('<i class="ti ti-checks me-1"></i>Configured');

            showAlertBox({
                title: "Success",
                message: "SendGrid configuration saved successfully",
                type: "info",
                sound: "/assets/sounds/success.mp3"
            });
        }
    } catch (err) {
        let errorMsg = "Failed to save SendGrid configuration";
        if (err.response && err.response.data && err.response.data.message) {
            errorMsg = err.response.data.message;
        }
        showAlertBox({ title: "Error", message: errorMsg, type: "danger", sound: "/assets/sounds/error.mp3" });
    } finally {
        btn.prop("disabled", false).html(originalHtml);
    }
}

async function saveSmtpSettings() {
    let host = $("#smtp_host").val().trim();
    let port = $("#smtp_port").val().trim();
    let username = $("#smtp_username").val().trim();
    let password = $("#smtp_password").val().trim();
    let encryption = $("#smtp_encryption").val();

    let hasError = false;

    if (!host) {
        $("#smtp_host").addClass("is-invalid");
        $("#error_msg_smtp_host").text("SMTP Host is required").show();
        hasError = true;
    } else {
        $("#smtp_host").removeClass("is-invalid");
        $("#error_msg_smtp_host").text("").hide();
    }

    if (!port || isNaN(port) || Number(port) < 1 || Number(port) > 65535) {
        $("#smtp_port").addClass("is-invalid");
        $("#error_msg_smtp_port").text("Valid Port number (1-65535) is required").show();
        hasError = true;
    } else {
        $("#smtp_port").removeClass("is-invalid");
        $("#error_msg_smtp_port").text("").hide();
    }

    if (!username) {
        $("#smtp_username").addClass("is-invalid");
        $("#error_msg_smtp_username").text("SMTP Username is required").show();
        hasError = true;
    } else {
        $("#smtp_username").removeClass("is-invalid");
        $("#error_msg_smtp_username").text("").hide();
    }

    if (hasError) {
        return;
    }

    let btn = $("#save_smtp_btn");
    let originalHtml = btn.html();
    btn.prop("disabled", true).html('<span class="spinner-border spinner-border-sm me-1" role="status"></span> Saving...');

    try {
        let formData = new FormData();
        formData.append("smtp_host", host);
        formData.append("smtp_port", port);
        formData.append("smtp_username", username);
        if (password) {
            formData.append("smtp_password", password);
        }
        formData.append("smtp_encryption", encryption);

        let response = await web.post('/settings/email-settings/smtp', formData);

        if (response.status === 200) {
            $("#modal_smtp").modal("hide");
            $("#smtp-status-badge")
                .removeClass("bg-light text-muted border fw-normal")
                .addClass("bg-success-transparent text-success fw-semibold")
                .html('<i class="ti ti-checks me-1"></i>Configured');

            showAlertBox({
                title: "Success",
                message: "SMTP configuration saved successfully",
                type: "info",
                sound: "/assets/sounds/success.mp3"
            });
        }
    } catch (err) {
        let errorMsg = "Failed to save SMTP configuration";
        if (err.response && err.response.data && err.response.data.message) {
            errorMsg = err.response.data.message;
        }
        showAlertBox({ title: "Error", message: errorMsg, type: "danger", sound: "/assets/sounds/error.mp3" });
    } finally {
        btn.prop("disabled", false).html(originalHtml);
    }
}

async function saveSesSettings() {
    let accessKeyId = $("#aws_access_key_id").val().trim();
    let secretAccessKey = $("#aws_secret_access_key").val().trim();
    let region = $("#aws_region").val();
    let configSet = $("#aws_config_set").val().trim();

    let hasError = false;

    if (!accessKeyId) {
        $("#aws_access_key_id").addClass("is-invalid");
        $("#error_msg_aws_access_key_id").text("AWS Access Key ID is required").show();
        hasError = true;
    } else {
        $("#aws_access_key_id").removeClass("is-invalid");
        $("#error_msg_aws_access_key_id").text("").hide();
    }

    if (!region) {
        $("#aws_region").addClass("is-invalid");
        $("#aws_region").next(".select2-container").find(".select2-selection").addClass("custom-is-invalid is-invalid");
        $("#error_msg_aws_region").text("AWS Region is required").show();
        hasError = true;
    } else {
        $("#aws_region").removeClass("is-invalid");
        $("#aws_region").next(".select2-container").find(".select2-selection").removeClass("custom-is-invalid is-invalid");
        $("#error_msg_aws_region").text("").hide();
    }

    if (hasError) {
        return;
    }

    let btn = $("#save_ses_btn");
    let originalHtml = btn.html();
    btn.prop("disabled", true).html('<span class="spinner-border spinner-border-sm me-1" role="status"></span> Saving...');

    try {
        let formData = new FormData();
        formData.append("aws_access_key_id", accessKeyId);
        if (secretAccessKey) {
            formData.append("aws_secret_access_key", secretAccessKey);
        }
        formData.append("aws_region", region);
        if (configSet) {
            formData.append("aws_config_set", configSet);
        }

        let response = await web.post('/settings/email-settings/ses', formData);

        if (response.status === 200) {
            $("#modal_ses").modal("hide");
            $("#ses-status-badge")
                .removeClass("bg-light text-muted border fw-normal")
                .addClass("bg-success-transparent text-success fw-semibold")
                .html('<i class="ti ti-checks me-1"></i>Configured');

            showAlertBox({
                title: "Success",
                message: "AWS SES configuration saved successfully",
                type: "info",
                sound: "/assets/sounds/success.mp3"
            });
        }
    } catch (err) {
        let errorMsg = "Failed to save AWS SES configuration";
        if (err.response && err.response.data && err.response.data.message) {
            errorMsg = err.response.data.message;
        }
        showAlertBox({ title: "Error", message: errorMsg, type: "danger", sound: "/assets/sounds/error.mp3" });
    } finally {
        btn.prop("disabled", false).html(originalHtml);
    }
}

async function sendTestEmail() {
    let recipient = $("#test_recipient_email").val().trim();
    let subject = $("#test_email_subject").val().trim();
    let message = $("#test_email_message").val().trim();
    let alertBox = $("#test_email_alert");

    alertBox.addClass("d-none").removeClass("alert-success alert-danger").text("");

    let emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!recipient || !emailRegex.test(recipient)) {
        $("#test_recipient_email").addClass("is-invalid");
        $("#error_msg_test_recipient_email").text("Valid recipient email is required").show();
        return;
    } else {
        $("#test_recipient_email").removeClass("is-invalid");
        $("#error_msg_test_recipient_email").text("").hide();
    }

    if (!subject) {
        $("#test_email_subject").addClass("is-invalid");
        $("#error_msg_test_email_subject").text("Subject is required").show();
        return;
    } else {
        $("#test_email_subject").removeClass("is-invalid");
        $("#error_msg_test_email_subject").text("").hide();
    }

    if (!message) {
        $("#test_email_message").addClass("is-invalid");
        $("#error_msg_test_email_message").text("Message content is required").show();
        return;
    } else {
        $("#test_email_message").removeClass("is-invalid");
        $("#error_msg_test_email_message").text("").hide();
    }

    let btn = $("#send_test_email_btn");
    let originalHtml = btn.html();
    btn.prop("disabled", true).html('<span class="spinner-border spinner-border-sm me-1" role="status"></span> Saving...');

    let formData = new FormData();
    formData.append("recipient_email", recipient);
    formData.append("subject", subject);
    formData.append("message", message);

    try {
        let response = await web.post('/settings/email-settings/test', formData);

        if (response.status === 200) {
            alertBox.removeClass("d-none").addClass("alert-success").text(response.data.message || "Test email sent successfully!");
            showAlertBox({
                title: "Test Email Sent",
                message: "Email successfully dispatched to " + recipient,
                type: "info",
                sound: "/assets/sounds/success.mp3"
            });
        } else {
            let errorMsg = response.data && response.data.message ? response.data.message : "Failed to send test email";
            alertBox.removeClass("d-none").addClass("alert-danger").text(errorMsg);
        }
    } catch (err) {
        let errorMsg = "An unexpected error occurred while sending test email";
        if (err.response && err.response.data && err.response.data.message) {
            errorMsg = err.response.data.message;
        }
        alertBox.removeClass("d-none").addClass("alert-danger").text(errorMsg);
    } finally {
        btn.prop("disabled", false).html(originalHtml);
    }
}

window.destroyPage = function () {
    if ($("#smtp_encryption").hasClass("select2-hidden-accessible")) {
        $("#smtp_encryption").select2("destroy");
    }
    if ($("#aws_region").hasClass("select2-hidden-accessible")) {
        $("#aws_region").select2("destroy");
    }
    $("#modal_smtp").off(".emailSelect2").off(".select2");
    $("#modal_ses").off(".emailSelect2").off(".select2");
    $("#aws_region").off(".emailRegion");
    $("#smtp_encryption").off(".emailSmtpEnc");
    $("#email-settings-edit-form").off(".emailSettings");
    $('input[name="mail_driver"]').off(".emailDriver");
    $(".toggle-password").off("click");
    $("#save_sendgrid_btn").off("click");
    $("#save_smtp_btn").off("click");
    $("#save_ses_btn").off("click");
    $("#send_test_email_btn").off("click");
};