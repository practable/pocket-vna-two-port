clear all

hs = sparameters('DUT 3 coupler S12.s2p');
hs2 = sparameters('DUT 2 s41 coupler.s2p');
% rfplot(hs,2,1,'-black')
% hold on

S12 = sparameters(hs)
S13 = sparameters(hs2)

% 
S12f =  unwrap(angle(rfparam(S12,2,1)))*180/pi
S13f =  unwrap(angle(rfparam(S13,2,1)))*180/pi
freq=(0.005:0.002496871:2)

figure(1)
plot(freq,S12f,freq,S13f)


Real_angle=(S12f-S13f)-180

figure(2)
plot(freq,Real_angle,'-r')
ylabel('Phase Imbalance in degrees')
xlabel('Frequency GHz')
xlim([0.4 2])
% ylim([-30 4])
grid on
hold on

% 

