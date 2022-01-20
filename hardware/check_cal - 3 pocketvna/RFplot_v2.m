clear all
clc
%%%%% Declaration in the table
C0=-4.3* 1e-15
C1= 0.1 * 1e-27
C2=-11.5*1e-36
C3=0.12*1e-45


%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%short calc

fid1 = fopen('short.s2p','r'); 
data = textscan(fid1,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data = cell2mat(data);

real1=formated_data(:,2);
imag1=formated_data(:,3);
Gs1short(:,1) = real1+j*imag1;

        f=formated_data(:,1);

%%%%% mm to m conversion of lenth
l=8.4973 * 10^-3

% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% f=1.5*10^9:0.006875000*10^9:7*10^9    

c= physconst('LightSpeed')
z0=50
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
lambda = c./f
Beta=(2*pi)./lambda
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% reflection of open

Ce= (C0) + (C1.*f)+(C2*f.^2)+(C3*f.^3)

Grefl_open = rot90 (((1-2*j.*pi.*f.*Ce.*z0)./(1+2.*pi.*Ce.*z0)).*exp(-2j.*Beta.*l))
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%




% figure(1)
% plot(f,real(Grefl_open),f,imag(Grefl_open))
% legend('real','imaginary')
% xlabel ('Frequency Hz')
% title('open')
% grid on


% figure(2)
% plot(f,Grefl_db_open)
% grid on
% ylabel ('dB')
% xlabel ('Frequency Hz')
% 
% ylim([-2 2])
% xlim([10 10*10^9])




%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%short calc

L=0

Grefl_short = rot90( ((j*2.*pi.*f.*L - z0)./(j.*2.*pi.*L + z0)).*exp(-2j.*Beta.*l) )

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%short calc




% figure(3)
% plot(f,real(Grefl_short),f,imag(Grefl_short))
% legend('real','imaginary')
% 
% % plot(f,Grefl_db_short)
% grid on
% title('short')
% ylabel ('dB')
% xlabel ('Frequency Hz')
% 





 
        
 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%open calc       
fid2 = fopen('open.s2p','r'); 
data2 = textscan(fid2,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data2 = cell2mat(data2);

real2=formated_data2(:,2);
imag2=formated_data2(:,3);
Gs2open(:,1) = real2+j*imag2;
        
 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%load      
        fid3 = fopen('load.s2p','r'); 
data3 = textscan(fid3,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data3 = cell2mat(data3);

real3=formated_data3(:,2);
imag3=formated_data3(:,3);
Gs3load(:,1) = real3+j*imag3;
                
 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%Gdutm     
                fid4 = fopen('DUTuncal.s2p','r'); 
data4 = textscan(fid4,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data4 = cell2mat(data4);

real4=formated_data4(:,2);
imag4=formated_data4(:,3);
Gdutm(:,1) = real4+j*imag4;

Grefl_load=zeros(length(f),1);  %Grefl_load=1 Cannot be one




for i=1:length(f)   %The index has to have two terms removed because it starts with 1 and 1 already
        A=[1,Grefl_short(i).*Gs1short(i),-Grefl_short(i);1,Grefl_open(i).*Gs2open(i),-Grefl_open(i);1,Grefl_load(i).*Gs3load(i),-Grefl_load(i)];
        B=[Gs1short(i);Gs2open(i);Gs3load(i)];
        X=linsolve(A,B); %Calculated error terms
        e00(i,1)=X(1,1); %Directivity
        e11(i,1)=X(2,1); %Source Match
        delta_e(i,1)=X(3,1);
end

e10e01=(e00.*e11)-delta_e; %Reflection Tracking


%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%Final part
R=(Gdutm-e00)./((Gdutm.*e11)-delta_e); %Real correct S11





 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%Gdutm     
fid5 = fopen('DUTcal.s2p','r');  
data5 = textscan(fid5,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data5 = cell2mat(data5);

real5=formated_data5(:,2);
imag5=formated_data5(:,3);
Gmes(:,1) = real5+j*imag5;

 
figure(1)
plot(f/10^9,db(R),'-red',f/10^9,db(Gmes),'-blue',f/10^9,db(Gdutm),'-black')
grid on
legend('Corrected from Process','Measured in reality','not calibrated')
xlabel('Frequency in GHz')
ylabel('Magnitude')

figure(2)
plot(f/10^9,unwrap(angle(R)),'-red',f/10^9,unwrap(angle(Gmes)),'-blue',f/10^9,unwrap(angle(Gdutm)),'-black')
grid on
legend('Corrected from Process','Measured in reality','not calibrated')
xlabel('Frequency in GHz')
ylabel('Angle')
%e10e01=e00*e11-delta_e
%Gdut= (1/e11)*((e00-Gdutm)/(e00-((e10e01)/e11)-Gdutm))


figure(3)
plot(f/10^9,db(e00),'-red',f/10^9,db(e11),'-blue',f/10^9,db(delta_e),'-black')
grid on
legend('Directivity','Source Match','Tracking?')
xlabel('Frequency in GHz')
ylabel('Angle')
