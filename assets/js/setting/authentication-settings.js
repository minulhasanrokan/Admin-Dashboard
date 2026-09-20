'use strict';

$(document).ready(async function () {

    bindAuthenticationSettingsEvents();
});

function bindAuthenticationSettingsEvents() {

    let form = $("#authentication-settings-form");

    form.off(".authSettings");

    form.on("focusout.authSettings", "input[type='text']", function () {
        validateField(this.id);
    });

    form.on("click.authSettings", "#authentication_setting_button", async function () {
        await saveAuthenticationSettings();
    });
}

function validateField(id) {

    switch (id) {
        case "verification_expired":
            form_validation({control: "verification_expired#number#mandatory", msg_text: "Verification Expired"});
            break;
    }
}

async function saveAuthenticationSettings() {

    let status = form_validation({
        control: "verification_expired#number#mandatory*login_type#radio#mandatory*otp_type#radio#mandatory",
        msg_text: "Verification Expired*Login Type*OTP Type"
    });

    if (status === false) {
        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
        return false;
    }

    let confirmed = await confirmDialog("Authentication Settings", "Are you sure you want to save the Authentication Settings?");

    if (!confirmed) {
        return;
    }

    let formData = new FormData();

    formData.append("allow_registration", $("#allow_registration").is(":checked"));
    formData.append("invite_only", $("#invite_only").is(":checked"));
    formData.append("verification_required", $("#verification_required").is(":checked"));
    formData.append("verification_expired", $("#verification_expired").val());
    formData.append("referral_system", $("#referral_system").is(":checked"));

    let loginType = $('input[name="login_type"]:checked').val() || "email";
    formData.append("login_type", loginType);

    formData.append("password_required", $("#password_required").is(":checked"));
    formData.append("otp_system", $("#otp_system").is(":checked"));

    let otpType = $('input[name="otp_type"]:checked').val() || "email";
    formData.append("otp_type", otpType);

    try {
        let response = await web.post('/settings/authentication-settings', formData);

        if (response.status === 200) {
            showAlertBox({title: "Success", message: "Authentication Settings Updated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
    } catch (err) {
    }
}

window.destroyPage = function () {
    $("#authentication-settings-form").off(".authSettings");
};
