function Display($http, $interval) {
    var display = this;
    this.isDispalyShown = false;

    this.key = function (key) {
        $http.post("display/key/" + key)
                .then(function (response) {
                    if (response.data !== true) {
                        alert(response.data);
                    }
                });
    }

    this.updateDisplay = function() {
        if (display.isDispalyShown) {
            $http.get("display/screen")
                .then(function (response) {
                    display.screen = response.data;
                });
        }
    };

    $interval(this.updateDisplay, 1000);
}