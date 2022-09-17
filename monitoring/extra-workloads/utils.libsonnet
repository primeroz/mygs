local kubecfg = import 'kubecfg.libsonnet';

local HelmToGVKName(parsedHelm) =
  local gvkName(accum, o) = accum {
    [o.apiVersion]+: {
      [o.kind]+: {
        assert !( o.metadata.name in super),
        [o.metadata.name]: o,
      },
    },
  };
  kubecfg.fold(
    gvkName,
    parsedHelm,
    {}
  );

{
  HelmToObjectsGVKName(parsedHelm):: HelmToGVKName(parsedHelm),
}
