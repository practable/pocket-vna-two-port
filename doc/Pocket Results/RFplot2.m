hs = sparameters('DUT 3 coupler S12.s2p');
rfplot(hs,2,1,'-red')
hold on

hs2 = sparameters('DUT 2 s41 coupler.s2p');
rfplot(hs2,2,1,'-blue')
hold on


hs3 = sparameters('DUT 4 coupler S13 isol.s2p');
rfplot(hs3,2,1,'-black')
hold on

legend('S12','S14','S13')
ylim([-53   0])