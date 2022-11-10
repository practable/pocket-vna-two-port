clear all
clc


%  %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% Calibrated using garbage VNA  and it is software  

       
fid_meas = fopen('dut_correct_data_pocket_system.s2p','r'); 
data_meas = textscan(fid_meas,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data_meas = cell2mat(data_meas);

real_meas=formated_data_meas(:,2);
imag_meas=formated_data_meas(:,3);
VNAcorS11(:,1) = real_meas+j*imag_meas;

real_meas2=formated_data_meas(:,4);
imag_meas2=formated_data_meas(:,5);
VNAcorS21(:,1) = real_meas2+j*imag_meas2;

real_meas3=formated_data_meas(:,6);
imag_meas3=formated_data_meas(:,7);
VNAcorS12(:,1) = real_meas3+j*imag_meas3;

real_meas4=formated_data_meas(:,8);
imag_meas4=formated_data_meas(:,9);
VNAcorS22(:,1) = real_meas4+j*imag_meas4;




%%%%% Declaration in the table
C0=-4.3* 1e-15
C1= 0.1 * 1e-27
C2=-11.5*1e-36
C3=0.12*1e-45


%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%short calc

fid1 = fopen('short_p1.s2p','r'); 
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
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%short calc

L=0

Grefl_short = rot90( ((j*2.*pi.*f.*L - z0)./(j.*2.*pi.*L + z0)).*exp(-2j.*Beta.*l) )
        
 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%open calc       
fid2 = fopen('open_p1.s2p','r'); 
data2 = textscan(fid2,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data2 = cell2mat(data2);

real2=formated_data2(:,2);
imag2=formated_data2(:,3);
Gs2open(:,1) = real2+j*imag2;
        
 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%load      
        fid3 = fopen('load_p1.s2p','r'); 
data8 = textscan(fid3,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data8 = cell2mat(data8);

real9=formated_data8(:,2);
imag9=formated_data8(:,3);
Gs3load(:,1) = real9+j*imag9;
              %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%  
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
  %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%Gdutm     UNCAL data of the
  %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%two port
fid4 = fopen('DUTuncal.s2p','r'); 
data4 = textscan(fid4,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data4 = cell2mat(data4);
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
real4=formated_data4(:,2);
imag4=formated_data4(:,3);
Gdutm(:,1) = real4+j*imag4;  
Gmes(:,1)= real4+j*imag4;  %%%%%%s11
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
realmescon=formated_data4(:,4); 
imagmescon=formated_data4(:,5);
Gmescon(:,1) = realmescon+j*imagmescon; %%%%%%s21
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
realGmescon_r=formated_data4(:,6);
imagGmescon_r=formated_data4(:,7);
Gmescon_r(:,1) = realGmescon_r+j*imagGmescon_r; %%%%%%s12
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
realGmes_r=formated_data4(:,8);
imagGmes_r=formated_data4(:,9);
Gmes_r(:,1) = realGmes_r+j*imagGmes_r;
Gdutm_r(:,1) = realGmes_r+j*imagGmes_r; %%%%%%s22
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
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
 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%Gdutm  
 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%Gdutm     
% fid5 = fopen('DUTcal.s2p','r');  
% data5 = textscan(fid5,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
% formated_data5 = cell2mat(data5);
% 
% real5=formated_data5(:,2);
% imag5=formated_data5(:,3);
% Gmes(:,1) = real5+j*imag5; %%%%measured calibrated
% 
%  
% figure(1)
% plot(f/10^9,db(R),'-red',f/10^9,db(Gmes),'-blue',f/10^9,db(Gdutm),'-black')
% grid on
% legend('Corrected from Process','Measured in reality','not calibrated')
% xlabel('Frequency in GHz')
% ylabel('Magnitude')
% 
% figure(2)
% plot(f/10^9,unwrap(angle(R)),'-red',f/10^9,unwrap(angle(Gmes)),'-blue',f/10^9,unwrap(angle(Gdutm)),'-black')
% grid on
% legend('Corrected from Process','Measured in reality','not calibrated')
% xlabel('Frequency in GHz')
% ylabel('Angle')
% %e10e01=e00*e11-delta_e
% %Gdut= (1/e11)*((e00-Gdutm)/(e00-((e10e01)/e11)-Gdutm))
% 

% figure(3)
% plot(f/10^9,db(e00),'-red',f/10^9,db(e11),'-blue',f/10^9,db(delta_e),'-black')
% grid on
% legend('Directivity','Source Match','Tracking?')
% xlabel('Frequency in GHz')
% ylabel('Angle')
 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
  %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
   %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
    %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
     %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
      %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%



      % %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% 
% %%%%%%%Connect ports with termination to get e30 e03 Leakage (good VNA no need
% %%%%%%%to connect termination
% 
%  filenames6='p1p2 uncal no bridge for s21.s2p';
%  data6 = read(rfdata.data,filenames6);
%     s_params6 = extract(data6,'S_PARAMETERS',50);

    %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%         e30(:,1) = s_params6(2,1,:); %Leakage
%         e03_r(:,1) = s_params6(1,2,:); % Reverse leakage FOR PORT 2

        e30(:,1) = zeros(length(f),1); %Leakage
        e03_r(:,1) = zeros(length(f),1); % Reverse leakage FOR PORT 2

 %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%



%%%%%%Connect ports together -> P1 + P2
fid5 = fopen('pocket_thru.s2p','r'); 
    data5 = textscan(fid5,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_data5 = cell2mat(data5);

realS11m=formated_data5(:,2);
imagS11m=formated_data5(:,3);

realS21m=formated_data5(:,4);
imagS21m=formated_data5(:,5);

realS12m=formated_data5(:,6);
imagS12m=formated_data5(:,7);

realS22m=formated_data5(:,8);
imagS22m=formated_data5(:,9);

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
S11M(:,1) = realS11m+j*imagS11m;
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
S21M(:,1) = realS21m+j*imagS21m;
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
S12M(:,1) = realS12m+j*imagS12m;
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
S22M(:,1) = realS22m+j*imagS22m;
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
e22=(S11M-e00)./(S11M.*e11-delta_e) %Port 2 match forward
e10e32=(S21M-e30).*(1-e11.*e22)   %transmission tracking
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%This is forward mode, now need
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%reverse

fidsp2 = fopen('short_p2.s2p','r'); 
datasp2  = textscan(fidsp2 ,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_datasp2 = cell2mat(datasp2);

realsp2=formated_datasp2(:,8);
imagsp2=formated_datasp2(:,9);
Gs4short(:,1)  = realsp2+j*imagsp2;




fidop2 = fopen('open_p2.s2p','r'); 
dataop2 = textscan(fidop2,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_dataop2 = cell2mat(dataop2);
realop2=formated_dataop2(:,8);
imagop2=formated_dataop2(:,9);
Gs5open(:,1)   = realop2+j*imagop2;


fidlp2 = fopen('load_p2.s2p','r'); 
datalp2 = textscan(fidlp2,'%f %f %f %f %f %f %f  %f %f %f %f','HeaderLines',4);
formated_datalp2 = cell2mat(datalp2);
reallp2=formated_datalp2(:,8);
imaglp2=formated_datalp2(:,9);
Gs6load(:,1)   = reallp2+j*imaglp2;



% 
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%Reverse mode

        for i=1:length(f)   
        C=[1,Grefl_short(i).*Gs4short(i),-Grefl_short(i);1,Grefl_open(i).*Gs5open(i),-Grefl_open(i);1,Grefl_load(i).*Gs6load(i),-Grefl_load(i)];
        D=[Gs4short(i);Gs5open(i);Gs6load(i)];
        Y=linsolve(C,D); %Calculated error terms
        e33_r(i,1)=Y(1,1); %Directivity
        e11_r(i,1)=Y(2,1); %Source Match
        delta_e_r(i,1)=Y(3,1);
end


Rs22=(Gdutm_r-e33_r)./((Gdutm_r.*e11_r)-delta_e_r);
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%  
e23e32_r=(e33_r.*e11_r)-delta_e_r; %Reflection Tracking reverse
e22_r=(S22M-e33_r)./(S22M.*e11_r-delta_e_r) %%%% Port 2 match reverse 
e23e01_r=(S12M-e03_r).*(1-e11_r.*e22_r) %%%% Transmission tracking reverse


%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%REVERSE MODE ENDS HERE
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%Measured Sparameters starts here
D=(1+((Gmes-e00)./e10e01).*e11).*(1+((Gmes_r-e33_r)./e23e32_r).*e22_r)...
    -((Gmescon-e30)./(e10e32)).*((Gmescon_r-e03_r)./(e23e01_r)).*e22.*e11_r    ;


% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
S11cor=((Gmes-e00)./e10e01).*(1+((Gmes_r-e33_r)./e23e32_r).*e22_r)-e22.*((Gmescon-e30)./(e10e32)).*((Gmescon_r-e03_r)./(e23e01_r))./D 
S22cor=((Gmes_r-e33_r)./e23e32_r).*(1+((Gmes-e00)./e10e01).*e11)-e11_r.*((Gmescon-e30)./(e10e32)).*((Gmescon_r-e03_r)./(e23e01_r))./D
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%




% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
S21cor=(((Gmescon-e30)./(e10e32)).*(1+(Gmes_r-e33_r)./(e23e32_r).*(e22_r-e22)))./D    ;
S12cor=(((Gmescon_r-e03_r)./(e23e01_r)).*(1+(Gmes-e00)./(e10e01).*(e11-e11_r)))./D    ;
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%



% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
figure(1) %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
plot(f/10^9,db(e00),'-red',f/10^9,db(e11),'-blue',f/10^9,db(delta_e),'-black')
grid on
title('Port 1 Error terms')
legend('Directivity','Source Match','Tracking?')
xlabel('Frequency in GHz')
ylabel('Error Term value')

figure(2) %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
plot(f/10^9,db(e33_r),'-red',f/10^9,db(e11_r),'-blue',f/10^9,db(delta_e_r),'-black')
grid on
title('Port 2 Error terms')
legend('Directivity','Source Match','Tracking?')
xlabel('Frequency in GHz')
ylabel('Error Term value')



figure(3) %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
plot(f/10^9,db(R),'-red',f/10^9,db(Rs22),'-blue',f/10^9,db(VNAcorS11),'--r',f/10^9,db(VNAcorS22),'--blue')
grid on
title('Corrected S11/S22')
legend('Corrected S11','Corrected S22','PocketS11','PocketS22')
xlabel('Frequency in GHz')
ylabel('dB')


figure(4) %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
plot(f/10^9,db(S21cor),'-red',f/10^9,db(S12cor),'-blue',f/10^9,db(VNAcorS21),'--r',f/10^9,db(VNAcorS12),'--blue')
grid on
title('Corrected S21/S12')
legend('Corrected S21','Corrected S12','PocketS21','PocketS12')
xlabel('Frequency in GHz')
ylabel('dB')

 