function Controller($http, $interval) {

    var controller = this;

    this.getInfo = function () {
        $http.get("controller/info")
            .then(function (response) {
                controller.info = response.data;
                var minutes = 1000 * 60;
                var hours = minutes * 60;
                var days = hours * 24;

                controller.info.Reminders.forEach(function (item) {
                    item.timeLeft = Math.round((new Date(item.Next) - Date.now()) / days);

                    item.reset = function () {
                        $http.post("command/resetReminder/" + item.Index)
                            .then(function (response) {
                                if (response.data !== true) {
                                    alert(response.data);
                                }

                                controller.getInfo();
                            });
                    }
                });

                controller.info.Reminders.IsOverdue = controller.info.Reminders
                    .some(function (reminder) { return reminder.IsOverdue });

                controller.info.Maintenance.forEach(function (item) {
                    item.toggle = function () {
                        $http.post("command/maintenance/" + item.Index, !item.IsActive)
                            .then(function (response) {
                                if (response.data !== true) {
                                    alert(response.data);
                                }

                                controller.getInfo();
                            });
                    }
                });

                controller.refreshAssoications();
            });
    }

    this.refreshAssoications = function () {
        var updateAssoications = function (mode) {
            if (mode.IsProbe && controller.probes) {
                mode.Icon = "icons/thermometer.svg";
                mode.Item = controller.probes.find(function (item) {
                    return item.Id === mode.Id;
                });

                if (mode.Item) {
                    mode.ValueString = mode.Item.ConvertedValue.toString() + mode.Item.Units;
                }
                return;
            }

            switch (mode.DeviceMode) {
                case "Lights":
                    mode.Icon = "icons/bulb.svg";
                    if (controller.lights) {
                        mode.Item = controller.lights.find(function (item) {
                            return item.Id === mode.Id;
                        });

                        if (mode.Item) {
                            mode.ValueString = mode.Item.Value.toString() + mode.Item.Units;
                        }
                    }
                    return;
                case "Timer":
                    mode.Icon = "icons/timer.svg";
                    if (controller.dosingPumps) {
                        mode.Item = controller.dosingPumps.find(function (item) {
                            return item.Id === mode.Id;
                        });

                        if (mode.Item) {
                            mode.Icon = "icons/dropper.svg";
                            mode.ValueString = mode.Item.Value;
                        }
                        else {
                            mode.Item = { DisplayName: "Timer " + mode.Port };
                        }
                    }
                    return;
                case "Water":
                    mode.Icon = "icons/water_level.svg";
                    if (controller.levelSensors) {
                        mode.Item = controller.levelSensors.find(function (item) {
                            return item.Id === mode.Id;
                        });

                        if (mode.Item) {
                            mode.ValueString = mode.Item.Value;
                        }
                    }
                    return;
                case "CurrentPump":
                    mode.Icon = "icons/wave.svg";
                    if (controller.pumps) {
                        mode.Item = controller.pumps.find(function (item) {
                            return item.Id === mode.Id;
                        });

                        if (mode.Item) {
                            mode.ValueString = mode.Item.Value.toString() + mode.Item.Units;
                        }
                    }

                    return;
                case "ProgrammableLogic":
                    mode.Icon = "icons/puzzle.svg";
                    if (controller.programmablelogic) {
                        mode.Item = controller.programmablelogic.find(function (item) {
                            return item.Index === parseInt(mode.Id);
                        });
                    }
                    return;
                case "Maintenance":
                    mode.Icon = "icons/wrench.svg";
                    if (controller.info) {
                        mode.Item = controller.info.Maintenance.find(function (item) {
                            return item.Index === mode.Port;
                        });

                        if (mode.Item) {
                            mode.ValueString = mode.Item.IsActive ? "On" : "Off";
                        }
                    }
                    return;
                case "ThunderStorm":
                    mode.Icon = "icons/thunder.svg";
                    mode.Item = { DisplayName: "Thunder Storm" };
                    return;
                case "Thunder":
                    mode.Icon = "icons/thunder.svg";
                    mode.Item = { DisplayName: "Thunder" };
                    return;
                case "Alarm":
                    mode.Icon = "icons/alarm.svg";
                    mode.Item = { DisplayName: "Alarm" };
                    return;
                case "WaterChange":
                    mode.Icon = "icons/water_change.svg";
                    mode.Item = { DisplayName: "Water Change" };
                    return;
                default:
            }
        }

        if (controller.sports) {
            controller.sports.forEach(function (item) {
                updateAssoications(item.Mode);
            });
        }

        if (controller.lports) {
            controller.lports.forEach(function (item) {
                updateAssoications(item.Mode);
            });
        }

        if (controller.programmablelogic) {
            controller.programmablelogic.forEach(function (item) {
                updateAssoications(item.Input1);
                updateAssoications(item.Input2);
                item.Input_list1 = [item.Input1];
                item.Input_list2 = [item.Input2];

                item.Input1.PlInvert = item.Function.Invert1;
                item.Input2.PlInvert = item.Function.Invert2;
            });
        }

        if (controller.levelSensors) {
            controller.levelSensors.forEach(function (item) {
                switch (item.OpertationMode) {
                    case "AutoTopOff":
                        item.Icon1 = "icons/water_level.svg";
                        break;
                    case "MinMaxControl":
                        item.Icon1 = "icons/water_level.svg";
                        break;
                    case "WaterChange":
                        item.Icon1 = "icons/water_change.svg";
                        break;
                    case "LeekageDetection":
                        item.Icon1 = "icons/alarm.svg";
                        break;
                    case "WaterChangeAndAutoTopOff":
                        item.Icon1 = "icons/water_level.svg";
                        item.Icon2 = "icons/water_change.svg";
                        break;
                    case "AutoTopOffWith2Sensors":
                        item.Icon1 = "icons/water_level.svg";
                        break;
                    case "ReturnPump":
                        item.Icon1 = "icons/wave.svg";
                        break;
                }
            });
        }
    }

    this.getData = function () {
        controller.getInfo();

        $http.get("controller/probe")
            .then(function (response) {
                controller.probes = response.data;
                controller.probes.IsAlarm = controller.probes.some(function (sensor) {
                    return sensor.AlarmState === "On";
                });
                controller.refreshAssoications();

                controller.probes.forEach(function (item) {
                    item.showGraph = function (range) {
                        controller.SelectedProbe = item;

                        controller.SelectedData = [];
                        if (controller.chart) {
                            controller.chart.update();
                        }

                        var path;

                        controller.SelectedGraph = range;
                        switch (range) {
                            case "day":
                                path = "data/log/";
                                break;
                            case "year":
                                path = "data/logYear/";
                                break;
                            default:
                                path = "data/logWeek/";
                                break;
                        }

                        $http.get(path + item.Id)
                            .then(function (response) {
                                response.data.forEach(function (dataPoint) {
                                    controller.SelectedData.push({parsedDate:Date.parse(dataPoint.time),  date: dataPoint.time,  x: Date.parse(dataPoint.time), y: dataPoint.value });
                                });

                                controller.SelectedData = controller.SelectedData.sort(function (a, b) { return a.parsedDate - b.parsedDate; });

                                if (!controller.chart) {
                                    var ctx = document.getElementById("lineChart");

                                    controller.chartConfig = {
                                        type: 'line',
                                        data: {
                                            datasets: [{
                                                data: controller.SelectedData,
                                                backgroundColor: [
                                                    'rgba(255, 99, 132, 0.2)'
                                                ],
                                                borderColor: [
                                                    'rgba(255,99,132,1)'
                                                ],
                                                borderWidth: 1
                                            }]
                                        },
                                        options: {
                                            title: {
                                                display: false
                                            },
                                            legend: { display: false },
                                            tooltips: { enabled: false },
                                            elements: { point: { radius: 0, hoverRadius: 0 } },
                                            maintainAspectRatio: false,
                                            scales: {
                                                xAxes: [{
                                                    type: "time",
                                                    display: true,
                                                    scaleLabel: {
                                                        display: true,
                                                        labelString: 'Time'
                                                    },
                                                    ticks: {
                                                        major: {
                                                            fontStyle: "bold",
                                                            fontColor: "#FF0000"
                                                        }
                                                    }
                                                }],
                                                yAxes: [{
                                                    display: true,
                                                    scaleLabel: {
                                                        display: true,
                                                        labelString: controller.SelectedProbe.Units
                                                    }
                                                }]
                                            }
                                        }
                                    };

                                    controller.chart = new Chart(ctx, controller.chartConfig);
                                }
                                else {
                                    controller.chartConfig.data.datasets[0].data = controller.SelectedData;
                                    controller.chartConfig.options.scales.yAxes[0].scaleLabel.labelString = controller.SelectedProbe.Units;
                                    controller.chart.update();
                                }
                            });
                    }
                });
            });

        $http.get("controller/levelsensor")
            .then(function (response) {
                controller.levelSensors = response.data;
                controller.levelSensors.IsAlarm = controller.levelSensors
                    .some(function (sensor) { return sensor.AlarmState === "On" });

                controller.levelSensors.forEach(function (item) {
                    item.clearAlarm = function () {
                        $http.post("command/clearlevelalarm/" + item.Id)
                            .then(function (response) {
                                if (response.data !== true) {
                                    alert(response.data);
                                }

                                controller.getData();
                            });
                    };

                    item.startWaterChange = function () {
                        $http.post("command/startwaterchange/" + item.Id)
                            .then(function (response) {
                                if (response.data !== true) {
                                    alert(response.data);
                                }

                                controller.getData();
                            });
                    };
                });

                controller.refreshAssoications();

            });

        $http.get("controller/sport")
            .then(function (response) {
                controller.sports = response.data;
                controller.sports.forEach(function (item) {
                    item.toggle = function () {
                        $http.post("command/SetSocket/" + item.Id, item.Value !== 'On')
                            .then(function (response) {
                                if (response.data !== true) {
                                    alert(response.data);
                                }

                                controller.getData();
                            });
                    };

                    item.showLogic = function () {

                        controller.refreshAssoications();

                        controller.SelectedLogic = [item.Mode];
                        controller.SelectedPort = item;

                        item.Mode
                    };
                });
                controller.refreshAssoications();
            });

        $http.get("controller/lport")
            .then(function (response) {
                controller.lports = response.data;
                controller.refreshAssoications();
            });

        $http.get("controller/digitalinput")
            .then(function (response) {
                controller.digitalInput = response.data;
            });

        $http.get("controller/pump")
            .then(function (response) {
                controller.pumps = response.data;
                controller.refreshAssoications();
            });

        $http.get("controller/programmablelogic")
            .then(function (response) {
                controller.programmablelogic = response.data;
                controller.refreshAssoications();
            });

        $http.get("controller/dosingpump")
            .then(function (response) {
                controller.dosingPumps = response.data;

                controller.dosingPumps.forEach(function (item) {
                    item.updateDousingValue = function (perDay, rate) {
                        $http.put("command/updatedousingvalue/" + item.Id, { PerDay: perDay, Rate: rate })
                            .then(function (response) {
                                if (response.data !== true) {
                                    alert(response.data);
                                }

                                controller.getData();
                            });
                    };
                });
                controller.refreshAssoications();
            });

        $http.get("controller/light")
            .then(function (response) {
                controller.lights = response.data;

                controller.lights.forEach(function (item) {
                    item.setLight = function (enable) {
                        $http.post("command/setlight/" + item.Id, enable)
                            .then(function (response) {
                                if (response.data !== true) {
                                    alert(response.data);
                                }

                                controller.getData();
                            });
                    };
                });
                controller.refreshAssoications();

            });
    }

    this.feedPause = function () {
        $http.post("command/feedpasue", true)
            .then(function (response) {
                if (response.data !== true) {
                    alert(response.data);
                }

                controller.getInfo();
            });
    }

    this.manualLights = function () {
        var enable = controller.info.OperationMode !== "ManualIllumination";
        $http.post("command/manuallights", enable)
            .then(function (response) {
                if (response.data !== true) {
                    alert(response.data);
                }

                controller.getInfo();
            });
    }

    this.manualSockets = function () {
        var enable = controller.info.OperationMode !== "ManualSockets";
        $http.post("command/manualSockets", enable)
            .then(function (response) {
                if (response.data !== true) {
                    alert(response.data);
                }

                controller.getInfo();
            });
    }

    this.thunderstorm = function () {
        $http.post("command/thunderstorm", 5)
            .then(function (response) {
                if (response.data !== true) {
                    alert(response.data);
                }

                controller.getInfo();
            });
    }

    this.refresh = function (enable) {
        $http.post("command/refresh", enable)
            .then(function (response) {
                if (response.data !== true) {
                    alert(response.data);
                }

                controller.getData();
            });
    }

    this.getData();
    $interval(this.getData, 10000);
}