# snippets

_Source: https://doc.photonengine.com/quantum/current/addons/bot-sdk/snippets_

# Snippets

Find here useful code snippets that can be used with Bot SDK.

## Compound Agents

Sometimes it might be useful to have more than one AI agent component operating in the same entity.

For example, maybe there is one HFSM which is responsible uniquely for moving the entity and another HFSM which is only responsible for attacking targets.

Even mixing different technologies in the same entity is possible: perhaps part of the work is done by with an HFSM and the other part if done with a BT.

When such setups are needed, it is also useful to provide more than one memory storage for the entity. In Bot SDK, it is common (but not mandatory) to use the AIBlackboard component.

Here are then code snippets which showcases how multiple components can be added to a single entity to create what we are calling here an entity with Compounds Agents:

### Creating the compound component

In the DSL (any `.qtn` file), add the component below.

Feel free to give them meaningful game-specific names to the components, or to add them in a generic-way to a list.

In this example, there is HFSM responsible only for Moving the entity and another one to perform Attacks (just to illustrate):

Qtn

```cs
component CompoundAgents
{
    HFSMAgent MovementHFSMAgent;
    HFSMAgent AttackHFSMAgent;
    AIBlackboardComponent MovementBlackboard;
    AIBlackboardComponent AttackBlackboard;
    AssetRefAIBlackboardInitializer MovementBBInitializer;
    AssetRefAIBlackboardInitializer AttackBBInitializer;
}

```

### Extending the AIContext

Create a partial implementation of `AIContextUser`. The context here is, just to exemplify: "given the updated entity, which of its HFSMs is being updated? Which of its Blackboards does that HFSM use?", etc:

C#

```csharp
namespace Quantum
{
  public unsafe partial struct AIContextUser
  {
    public readonly AIBlackboardComponent* Blackboard;
    public readonly HFSMAgent* HFSMAgent;
    public AIContextUser(AIBlackboardComponent* blackboard, HFSMAgent* hfsmAgent)
    {
      Blackboard = blackboard;
      HFSMAgent = hfsmAgent;
    }
  }
}

```

### Initialising the Agents and Blackboards

This is exemplified with a ISignalOnComponentAdded, but it is also up to user preference:

C#

```csharp
public void OnAdded(Frame frame, EntityRef entity, CompoundAgents* compoundAgents)
{
    // Initialise Agents
    HFSMRoot hfsmRoot = frame.FindAsset<HFSMRoot>(compoundAgents->MovementHFSMAgent.Data.Root.Id);
    HFSMManager.Init(frame, &compoundAgents->MovementHFSMAgent.Data, entity, hfsmRoot);
    hfsmRoot = frame.FindAsset<HFSMRoot>(compoundAgents->AttackHFSMAgent.Data.Root.Id);
    HFSMManager.Init(frame, &compoundAgents->AttackHFSMAgent.Data, entity, hfsmRoot);
    // Initialise Blackboards
    AIBlackboardInitializer initializer = frame.FindAsset<AIBlackboardInitializer>(compoundAgents->MovementBBInitializer.Id);
    AIBlackboardInitializer.InitializeBlackboard(frame, &compoundAgents->MovementBlackboard, initializer);
    initializer = frame.FindAsset<AIBlackboardInitializer>(compoundAgents->AttackBBInitializer.Id);
    AIBlackboardInitializer.InitializeBlackboard(frame, &compoundAgents->AttackBlackboard, initializer);
}

```

### Updating all Agents

Create the AI Context and fill it with each HFSM and Blackboard, then perform the Update:

C#

```csharp
public override void Update(Frame frame, ref Filter filter)
{
    HFSMData* movementData = &filter.CompoundAgents->MovementHFSMAgent.Data;
    HFSMData* rotationData = &filter.CompoundAgents->RotationHFSMAgent.Data;
    HFSMData* attackData = &filter.CompoundAgents->AttackHFSMAgent.Data;
    AIContext aiContext = new AIContext();
    AIContextUser movementContext = new AIContextUser(&filter.CompoundAgents->MovementHFSMAgent, &filter.CompoundAgents->MovementBlackboard);
    aiContext.UserData = &userData;
    HFSMManager.Update(frame, frame.DeltaTime, movementData, filter.EntityRef, ref aiContext);
    AIContextUser rotationContext = new AIContextUser(&filter.CompoundAgents->RotationHFSMAgent, &filter.CompoundAgents->RotationBlackboard);
    aiContext.UserData = &userData;
    HFSMManager.Update(frame, frame.DeltaTime, rotationData, filter.EntityRef, ref aiContext);
    AIContextUser attackContext = new AIContextUser(&filter.CompoundAgents->AttackHFSMAgent, &filter.CompoundAgents->AttackBlackboard);
    aiContext.UserData = &userData;
    HFSMManager.Update(frame, frame.DeltaTime, attackData, filter.EntityRef, ref aiContext);
}

```

### Getting the HFSM-specific blackboard from the AIContext object

Here is a snippet of an Action which writes to the context-specific Blackboard. Notice it is agnostic to which one it is, so the code is decoupled, there is no need to add lots of boilerplate code in the Actions to figure out which Blackboard should be used.

There is an _extension method_ which converts to the user context type:

C#

```csharp
namespace Quantum
{
    [System.Serializable]
    public unsafe class WriteBlackboardCompound : AIAction
    {
        public int Value;
        public override void Update(Frame frame, EntityRef entity, ref AIContext aiContext)
        {
            aiContext.UserData().Blackboard->Set(frame, "TestInteger", Value);
        }
    }
}

```

Back to top

- [Compound Agents](#compound-agents)
  - [Creating the compound component](#creating-the-compound-component)
  - [Extending the AIContext](#extending-the-aicontext)
  - [Initialising the Agents and Blackboards](#initialising-the-agents-and-blackboards)
  - [Updating all Agents](#updating-all-agents)
  - [Getting the HFSM-specific blackboard from the AIContext object](#getting-the-hfsm-specific-blackboard-from-the-aicontext-object)