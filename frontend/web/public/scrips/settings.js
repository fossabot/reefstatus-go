function Settings($http) {
    var settings = this;
    this.getData = function() {
        $http.get("settings/connection")
            .then(function(response) {
                settings.connection = response.data;
            });

        $http.get("settings/logging")
            .then(function(response) {
                settings.logging = response.data;
            });

        //$http.get("settings/mail")
        //    .then(function(response) {
        //        settings.mail = response.data;
        //    });
    };

    this.setData = function () {
        $http.put("settings/connection", this.connection)
            .then(function (response) {
                if (response.data !== true) {
                    alert(response.data);
                }
            });

        $http.put("settings/logging", this.logging)
            .then(function (response) {
                if (response.data !== true) {
                    alert(response.data);
                }
            });

        //$http.put("settings/mail", this.mail)
        //    .then(function (response) {
        //        if (response.data !== true) {
        //            alert(response.data);
        //        }
        //    });
    };

    this.getData();
}