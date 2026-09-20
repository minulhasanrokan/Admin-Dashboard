'use strict';

$(document).ready(async function () {

    bindUserRoleSettingsEvents();
});

function bindUserRoleSettingsEvents() {

    let form = $("#user-role-settings-edit-form");

    form.off(".userRoleSettings");

    form.on("focusout.userRoleSettings", "input", function () {validateField(this.id);});

    form.on("click.userRoleSettings", "#user_role_setting_button", async function () {

        await saveUserRoleSettings();
    });
}

function validateField(id) {

    switch (id) {
        case "password_min_length":
            form_validation({control: "password_min_length#number#mandatory", msg_text: "Password Min Length"});
            break;
        case "password_max_length":
            form_validation({control: "password_max_length#number#mandatory", msg_text: "Password Max Length"});
            break;
        case "session_time":
            form_validation({control: "session_time#number#mandatory", msg_text: "Session Time"});
            break;
        case "session_check_time":
            form_validation({control: "session_check_time#number#mandatory", msg_text: "Session Check Time"});
            break;
        case "password_expire_days":
            form_validation({control: "password_expire_days#number#mandatory", msg_text: "Password Expire Days"});
            break;
    }
}

async function saveUserRoleSettings() {

    let status = form_validation({control: "password_min_length#number#mandatory*password_max_length#number#mandatory*session_time#number#mandatory*session_check_time#number#mandatory*password_expire_days#number#mandatory", msg_text: "Password Min Length*Password Max Length*Session Time*Session Check Time*Password Expire Days"});

    if (status === false) {

        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
        return false;
    }

    let password_min_length = $("#password_min_length").val();
    let password_max_length = $("#password_max_length").val();

    if(parseFloat(password_max_length) < parseFloat(password_min_length)) {
        
        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
        showErrorMsg("password_min_length", "Password Max Length Can Not Less than Password Min Length");

        return;
    }

    let confirmed = await confirmDialog("User Role Settings", "Are you sure you want to save the User Role Settings?");

    if (!confirmed) {
        return;
    }

    let formData = new FormData();

    formData.append("password_min_length", password_min_length);
    formData.append("password_max_length", password_max_length);
    formData.append("session_time", $("#session_time").val());
    formData.append("session_check_time", $("#session_check_time").val());
    formData.append("password_expire_days", $("#password_expire_days").val());

    let response = await web.post('/settings/user-roles-settings', formData);

    if (response.status === 200) {

        showAlertBox({title: "Success", message: "User Role Settings Updated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
    }
}

window.destroyPage = function () {

    $("#user-role-settings-edit-form").off(
        ".userRoleSettings"
    );
};