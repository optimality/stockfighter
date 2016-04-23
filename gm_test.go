package stockfighter

import "testing"

func TestGM(t *testing.T) {
	client := MakeStockfighterClient(t)

	instance, err := client.StartLevel("first_steps")
	if err != nil {
		t.Fatalf("Error starting level: %v\n", err)
	}

	instance, err = client.ResumeInstance(instance.InstanceID)
	if err != nil {
		t.Fatalf("Error resuming level: %v\n", err)
	}

	instance, err = client.RestartInstance(instance.InstanceID)
	if err != nil {
		t.Fatalf("Error restarting level: %v\n", err)
	}

	_, err = client.StopInstance(instance.InstanceID)
	if err != nil {
		t.Fatalf("Error stopping level: %v\n", err)
	}
}
