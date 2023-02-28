% Author: Alexander Don
% University of Edinburgh (s2438345@ed.ac.uk)
% Remote Lab - Prototype test
% Date: 10-12-2022

%(Last Ammended 11-01-2023)

clc
clear all
close all

%Adding the measurement folders to the workspace
addpath("PocketVNA\")
addpath("Copper Mountain VNA\")

% Extracting the S-Parameter Data

%PocketVNA Measurements
PVNA_DUT1=sparameters('PVNA_DUT_1.s2p');
PVNA_DUT2=sparameters('PVNA_DUT_2.s2p');
PVNA_DUT3=sparameters('PVNA_DUT_3.s2p');
PVNA_DUT4=sparameters('PVNA_DUT_4.s2p');

% CopperMountain VNA (Professional VNA)
CMVNA_DUT1=sparameters('DUT_1.s2p');
CMVNA_DUT2=sparameters('DUT_2.s2p');
CMVNA_DUT3=sparameters('DUT_3.s2p');
CMVNA_DUT4=sparameters('DUT_4.s2p');

% "Ideal" Device reference S-parameters ( i.e the device measured soley
% using a VNA (not using the system)

DUT_1_REF= sparameters('DUT_1_REF.s2p');
DUT_2_REF= sparameters('DUT_2_REF.s2p');
DUT_3_REF= sparameters('DUT_3_REF.s2p');
DUT_4_REF= sparameters('DUT_4_REF.s2p');

%% DUT 1
figure
sgtitle('DUT 1 (Stub Filter)')

subplot(2,2,1)
rfplot(PVNA_DUT1,1,1,'db');
hold on
rfplot(CMVNA_DUT1,1,1,'db');
hold on
rfplot(DUT_1_REF,1,1,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 0.8);
title('S11','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,2)
rfplot(PVNA_DUT1,1,2,'db');
hold on
rfplot(CMVNA_DUT1,1,2,'db');
hold on
rfplot(DUT_1_REF,1,2,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S12','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,3)
rfplot(PVNA_DUT1,2,2,'db');
hold on
rfplot(CMVNA_DUT1,2,2,'db');
hold on
rfplot(DUT_1_REF,2,2,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S22','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,4)
rfplot(PVNA_DUT1,2,1,'db');
hold on
rfplot(CMVNA_DUT1,2,1,'db');
hold on
rfplot(DUT_1_REF,2,1,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S21','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

%% DUT 2
figure
sgtitle('DUT 2 (Ratrace Coupler) (Ports 1 & 2')

subplot(2,2,1)
rfplot(PVNA_DUT2,1,1,'db');
hold on
rfplot(CMVNA_DUT2,1,1,'db');
hold on
rfplot(DUT_2_REF,1,1,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 0.8);
title('S11','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor


subplot(2,2,2)
rfplot(PVNA_DUT2,1,2,'db');
hold on
rfplot(CMVNA_DUT2,1,2,'db');
hold on
rfplot(DUT_2_REF,1,2,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S12','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,3)
rfplot(PVNA_DUT2,2,2,'db');
hold on
rfplot(CMVNA_DUT2,2,2,'db');
hold on
rfplot(DUT_2_REF,2,2,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S22','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,4)
rfplot(PVNA_DUT2,2,1,'db');
hold on
rfplot(CMVNA_DUT2,2,1,'db');
hold on
rfplot(DUT_2_REF,2,1,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S21','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

%% DUT 3
figure
sgtitle('DUT 3 (Ratrace Coupler) (Ports 1 & 4')

subplot(2,2,1)
rfplot(PVNA_DUT3,1,1,'db');
hold on
rfplot(CMVNA_DUT3,1,1,'db');
hold on
rfplot(DUT_3_REF,1,1,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 0.8);
title('S11','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor


subplot(2,2,2)
rfplot(PVNA_DUT3,1,2,'db');
hold on
rfplot(CMVNA_DUT3,1,2,'db');
hold on
rfplot(DUT_3_REF,1,2,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S12','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,3)
rfplot(PVNA_DUT3,2,2,'db');
hold on
rfplot(CMVNA_DUT3,2,2,'db');
hold on
rfplot(DUT_3_REF,2,2,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S22','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,4)
rfplot(PVNA_DUT3,2,1,'db');
hold on
rfplot(CMVNA_DUT3,2,1,'db');
hold on
rfplot(DUT_3_REF,2,1,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S21','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

%% DUT 4
figure
sgtitle('DUT 4 (Ratrace Coupler) (Ports 1 & 3)')

subplot(2,2,1)
rfplot(PVNA_DUT4,1,1,'db');
hold on
rfplot(CMVNA_DUT4,1,1,'db');
hold on
rfplot(DUT_4_REF,1,1,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 0.8);
title('S11','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor


subplot(2,2,2)
rfplot(PVNA_DUT4,1,2,'db');
hold on
rfplot(CMVNA_DUT4,1,2,'db');
hold on
rfplot(DUT_4_REF,1,2,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S12','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,3)
rfplot(PVNA_DUT4,2,2,'db');
hold on
rfplot(CMVNA_DUT4,2,2,'db');
hold on
rfplot(DUT_4_REF,2,2,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S22','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor

subplot(2,2,4)
rfplot(PVNA_DUT4,2,1,'db');
hold on
rfplot(CMVNA_DUT4,2,1,'db');
hold on
rfplot(DUT_4_REF,2,1,'db');
set(findobj(gca, 'Type', 'Line'), 'LineWidth', 1);
title('S21','Interpreter','latex','fontsize',13)
ylim([-50,50])
yticks(-50:10:50)
legend('PVNA','CMVNA','Ideal Reference','Interpreter','Latex')
grid minor