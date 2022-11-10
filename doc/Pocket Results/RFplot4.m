hs = sparameters('DUT 3 coupler S12.s2p');
rfplot(hs,2,1,'angle','-red')
hold on

hs2 = sparameters('DUT 2 s41 coupler.s2p');
rfplot(hs2,2,1,'angle','-blue')
hold on


legend('S12','S14','S13')
% ylim([-53   0])