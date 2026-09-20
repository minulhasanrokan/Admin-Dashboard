'use strict';

$(document).ready(async function () {

    let country = $("#country_code").val();

    await loadDropDown('countries', 'code,name', 'country', 'country_container', 'Select Country', 0, country, 0, "8", "0", "0", 0, 0);

    let state = $("#state_code").val();

    await loadDropDown('states', 'code,name', 'state', 'state_container', 'Select State', 0, state, 0, "9", "country_code", country, 0, 0);

    let city = $("#city_code").val();

    await loadDropDown('cities', 'code,name', 'city', 'city_container', 'Select City', 0, city, 0, "10", "state_code", state, 0, 0);

    bindBusinessSettingsEvents();
});

function bindBusinessSettingsEvents() {

    let form = $("#business-settings-form");

    form.off(".businessSettings");

    form.on("focusout.businessSettings", "input, textarea, select", function () {validateField(this.id);});
    
    form.on("select2:close.businessSettings", "select", function () {validateField(this.id);});

    form.on("change.businessSettings", "#country", async function () {

        let country = $("#country").val();

        await loadDropDown('states', 'code,name', 'state', 'state_container', 'Select State', 0, "", 0, "9", "country_code", country, 0, 0);
    });

    form.on("change.businessSettings", "#state", async function () {

        let state = $("#state").val();

        await loadDropDown('cities', 'code,name', 'city', 'city_container', 'Select City', 0, "", 0, "10", "state_code", state, 0, 0);
    });

    form.on("click.businessSettings", "#business_setting_button", async function () {

        await saveBusinessSettings();
    });
}

function validateField(id) {

    switch (id) {
        case "name":
            form_validation({control: "name#text#mandatory", msg_text: "Name"});
            break;
        case "title":
            form_validation({control: "title#text#mandatory", msg_text: "Title"});
            break;
        case "slogan":
            form_validation({control: "slogan#text#mandatory", msg_text: "Slogan"});
            break;
        case "email":
            form_validation({control: "email#email#mandatory", msg_text: "Email"});
            break;
        case "mobile":
            form_validation({control: "mobile#mobile#mandatory", msg_text: "Mobile Number"});
            break;
        case "fax":
            form_validation({control: "fax#text#mandatory", msg_text: "Fax"});
            break;
        case "address":
            form_validation({control: "address#text#mandatory", msg_text: "Address"});
            break;
        case "country":
            form_validation({control: "country#select#mandatory", msg_text: "Country"});
            break;
        case "state":
            form_validation({control: "state#select#mandatory", msg_text: "State"});
            break;
        case "city":
            form_validation({control: "city#select#mandatory", msg_text: "City"});
            break;
        case "post_code":
            form_validation({control: "post_code#number#mandatory", msg_text: "Post Code"});
            break;
        case "details":
            form_validation({control: "details#text#nonmandatory", msg_text: "Details"});
            break;
        case "version":
            form_validation({control: "version#text#mandatory", msg_text: "Version"});
            break;
        case "copy_right":
            form_validation({control: "copy_right#text#mandatory", msg_text: "Copy Right"});
            break;
    }
}

async function saveBusinessSettings() {

    let status = form_validation({control: "name#text#mandatory*title#text#mandatory*slogan#text#mandatory*email#email#mandatory*mobile#mobile#mandatory*fax#text#mandatory*address#text#mandatory*country#select#mandatory*state#select#mandatory*city#select#mandatory*post_code#number#mandatory*details#text#nonmandatory*version#text#mandatory*copy_right#text#mandatory", msg_text: "Name*Title*Slogan*Email*Mobile Number*Fax*Address*Country*State*City*Post Code*Details*Version*Copy Right"});

    if (status === false) {

        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
        return false;
    }

    let confirmed = await confirmDialog("Business Settings", "Are you sure you want to save the Business Settings?");

    if (!confirmed) {
        return;
    }
    
    let formData = new FormData();

    formData.append("name", $("#name").val());
    formData.append("title", $("#title").val());
    formData.append("slogan", $("#slogan").val());
    formData.append("email", $("#email").val());
    formData.append("mobile", $("#mobile").val());
    formData.append("fax", $("#fax").val());
    formData.append("address", $("#address").val());
    formData.append("country", $("#country").val());
    formData.append("state", $("#state").val());
    formData.append("city", $("#city").val());
    formData.append("post_code", $("#post_code").val());
    formData.append("details", $("#details").val());
    formData.append("version", $("#version").val());
    formData.append("copy_right", $("#copy_right").val());

    let response = await web.post('/settings/business-settings', formData);

    if (response.status === 200) {

        showAlertBox({title: "Success", message: "Business Settings Updated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
    }
}

window.destroyPage = function () {

    $("#country, #state, #city").each(
        function () {

            if (
                $(this).hasClass(
                    "select2-hidden-accessible"
                )
            ) {

                $(this).select2(
                    "destroy"
                );
            }
        }
    );

    $("#business-settings-form").off(
        ".businessSettings"
    );
};