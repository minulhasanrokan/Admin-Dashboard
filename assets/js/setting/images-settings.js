'use strict';

$(document).ready(async function () {

    bindImagesSettingsEvents();
});

function bindImagesSettingsEvents() {

    let form = $("#images-settings-form");

    form.off(".imagesSettings");

    form.on("change.imagesSettings", ".image-sign",
        function () {
            showImagePreview(this);
        }
    );

    form.on("click.imagesSettings", ".image-cancel",
        function (e) {
            e.preventDefault();
            cancelImagePreview(this);
        }
    );

    form.on("click.imagesSettings", "#images_setting_button",
        async function () {
            await saveImagesSettings();
        }
    );

    form.on("change.imagesSettings", ".image-sign", function () {validateField(this.id);});

    form.on("click.imagesSettings", "#images_setting_button", async function () {await saveImageSettings();});
}

async function saveImageSettings() {

    let status = form_validation({control: "logo#file#nonmandatory*logo_white#file#nonmandatory*logo_small#file#nonmandatory*logo_dark_small#file#nonmandatory*icon#file#nonmandatory*apple_icon#file#nonmandatory*login_bg#file#nonmandatory*forgot_password_bg#file#nonmandatory*registration_bg#file#nonmandatory*verification_bg#file#nonmandatory", msg_text: "White Logo*Dark Logo*White Mini Logo*Dark Mini Logo*Favicon*Apple Icon*Auth Background*Forgot Password Background*Registration Background*Verification Background"});

    if (status === false) {

        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
        return false;
    }

    let confirmed = await confirmDialog("Images Settings", "Are you sure you want to save the Images Settings?");

    if (!confirmed) {
        return;
    }

    let formData = new FormData();

    let logo = document.getElementById("logo");
    if (logo && logo.files.length > 0) {
        formData.append("logo", logo.files[0]);
    }

    let logoWhite = document.getElementById("logo_white");
    if (logoWhite && logoWhite.files.length > 0) {
        formData.append("logo_white", logoWhite.files[0]);
    }

    let logoSmall = document.getElementById("logo_small");
    if (logoSmall && logoSmall.files.length > 0) {
        formData.append("logo_small", logoSmall.files[0]);
    }

    let logoDarkSmall = document.getElementById("logo_dark_small");
    if (logoDarkSmall && logoDarkSmall.files.length > 0) {
        formData.append("logo_dark_small", logoDarkSmall.files[0]);
    }

    let icon = document.getElementById("icon");
    if (icon && icon.files.length > 0) {
        formData.append("icon", icon.files[0]);
    }

    let appleIcon = document.getElementById("apple_icon");
    if (appleIcon && appleIcon.files.length > 0) {
        formData.append("apple_icon", appleIcon.files[0]);
    }

    let loginBg = document.getElementById("login_bg");
    if (loginBg && loginBg.files.length > 0) {
        formData.append("login_bg", loginBg.files[0]);
    }

    let forgotPasswordBg = document.getElementById("forgot_password_bg");
    if (forgotPasswordBg && forgotPasswordBg.files.length > 0) {
        formData.append("forgot_password_bg", forgotPasswordBg.files[0]);
    }

    let registrationBg = document.getElementById("registration_bg");
    if (registrationBg && registrationBg.files.length > 0) {
        formData.append("registration_bg", registrationBg.files[0]);
    }

    let verificationBg = document.getElementById("verification_bg");
    if (verificationBg && verificationBg.files.length > 0) {
        formData.append("verification_bg", verificationBg.files[0]);
    }

    let response = await web.post('/settings/images-settings', formData);

    if (response.status === 200) {

        showAlertBox({title: "Success", message: "Images Settings Updated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
    }
}

function validateField(id) {

    switch (id) {
        case "logo":
            form_validation({control: "logo#file#nonmandatory", msg_text: "White Logo"});
            break;
        case "logo_white":
            form_validation({control: "logo_white#file#nonmandatory", msg_text: "Dark Logo"});
            break;
        case "logo_small":
            form_validation({control: "logo_small#file#nonmandatory", msg_text: "White Mini Logo"});
            break;
        case "logo_dark_small":
            form_validation({control: "logo_dark_small#file#nonmandatory", msg_text: "Dark Mini Logo"});
            break;
        case "icon":
            form_validation({control: "icon#file#nonmandatory", msg_text: "Favicon"});
            break;
        case "apple_icon":
            form_validation({control: "apple_icon#file#nonmandatory", msg_text: "Apple Icon"});
            break;
        case "login_bg":
            form_validation({control: "login_bg#file#nonmandatory", msg_text: "Auth Background"});
            break;
        case "forgot_password_bg":
            form_validation({control: "forgot_password_bg#file#nonmandatory", msg_text: "Forgot Password Background"});
            break;
        case "registration_bg":
            form_validation({control: "registration_bg#file#nonmandatory", msg_text: "Registration Background"});
            break;
        case "verification_bg":
            form_validation({control: "verification_bg#file#nonmandatory", msg_text: "Verification Background"});
            break;
    }
}

function showImagePreview(input) {

    let $input = $(input);

    if (!$input.length) {
        return;
    }

    let inputId = $input.attr("id");

    if (!inputId) {
        return;
    }

    let $image = $("#" + inputId + "_image");

    if (!$image.length) {
        return;
    }

    let file = input.files && input.files[0];

    if (!file) {
        return;
    }

    if (!file.type.startsWith("image/")) {

        showAlertBox({title: "Error", message: "Please select a valid image file.", type: "danger", sound: "/assets/sounds/error.mp3"});

        $input.val("");

        return;
    }

    if (!$image.attr("data-original-src")) {

        $image.attr("data-original-src", $image.attr("src") || "");
    }

    let reader = new FileReader();

    reader.onload = function (e) {

        $image.attr("src", e.target.result);
    };

    reader.readAsDataURL(file);
}

function cancelImagePreview(cancelButton) {

    let $cancel = $(cancelButton);

    if (!$cancel.length) {
        return;
    }

    let inputId = $cancel.attr("data-input");

    if (!inputId) {
        return;
    }

    let $input = $("#" + inputId);
    let $image = $("#" + inputId + "_image");

    $input.val("");

    let originalSrc = $image.attr("data-original-src");

    if (originalSrc !== undefined) {

        $image.attr("src", originalSrc);
    }
}

window.destroyPage = function () {

    $("#images-settings-form").off(
        ".imagesSettings"
    );
};